package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/cxrdevelop/optimization_engine/api_server/internal/models"
	"github.com/cxrdevelop/optimization_engine/api_server/internal/optimization"
	"github.com/cxrdevelop/optimization_engine/pkg/compressor"
	"github.com/cxrdevelop/optimization_engine/pkg/environment"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/cxrdevelop/optimization_engine/pkg/storage"
)

const (
	ErrMsgEmptyRequest = "empty request body"
	ErrMsgJsonParse    = "error parsing json body"
	ErMsgJsonValidate  = "error validating json body"
	ErrMsgInternal     = "internal error"
	ErrMsgDownload     = "failed to download files"
	ErrMsgScript       = "script error"
	ErrMsgUpload       = "failed to upload the result"

	envPrefix     = "tmp_uploaded_"
	inputFileName = "input_files"

	maxMemory = 10 << 20 // 10 MB buffer
)

type UploadHandler struct {
	storage storage.Storage
	client  *optimization.Client
	log     *logger.Logger
}

func NewUploadHandler(storage storage.Storage, client *optimization.Client, log *logger.Logger) *UploadHandler {
	return &UploadHandler{
		storage: storage,
		client:  client,
		log:     log,
	}
}

func (h *UploadHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	// Create a temporary directory
	env := environment.New(os.TempDir(), envPrefix)
	if err := env.CreateTempDir(); err != nil {
		h.log.Errorf("error creating tempdir")
		writeResponse(writer, models.NewErrorResponse(err.Error()), http.StatusInternalServerError, h.log)
	}
	defer func() {
		if err := env.CleanUp(); err != nil {
			h.log.Warnf("error cleaning up temporary directory %s, error: %s", env.Dir(), err)
		}
	}()

	// Upload files from the UI
	h.log.Debugf("Start uploading files")
	filenames, err := h.serveFileUpload(env.Dir(), r)
	if err != nil {
		writeResponse(writer, models.NewErrorResponse(ErrMsgInternal), http.StatusBadRequest, h.log)
		return
	}

	// Create zip archive from user provided files
	h.log.Debugf("Add files to zip archive: '%s'", filenames)
	absPathToArch := path.Join(env.Dir(), (environment.Filename)(inputFileName).WithUnixSuffix())
	archFilesPath, err := compressor.Compress(context.Background(), absPathToArch, env.Dir(), filenames...)
	if err != nil {
		writeResponse(writer, models.NewErrorResponse(ErrMsgInternal), http.StatusBadRequest, h.log)
		return
	}

	// Upload files to the bucket
	filename := filepath.Base(archFilesPath)
	if _, err := h.storage.UploadFiles(env.Dir(), filename); err != nil {
		writeResponse(writer, models.NewErrorResponse("error uploading archive to the bucket"), http.StatusBadRequest, h.log)
		return
	}

	// POST a request to the optimization service
	h.log.Debugf("Post a request to the optimization service, filename: %s", filename)
	resp, err := h.client.PostOptimize(filename)
	if err != nil {
		writeResponse(writer, models.NewErrorResponse("script execution error"), http.StatusBadRequest, h.log)
		return
	}

	// Write response
	writeResponse(writer, models.NewUploadResponse(resp.Location, resp.Filepath, resp.ETag), http.StatusOK, h.log)
}

func (h *UploadHandler) serveFileUpload(workDir string, r *http.Request) ([]string, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, fmt.Errorf("error parsing file: %s", err)
	}

	formdata := r.MultipartForm.File["file"]
	filenames := make([]string, 0, len(formdata))
	for _, header := range formdata {
		filename, err := h.uploadFile(workDir, header)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func (h *UploadHandler) uploadFile(workDir string, header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("error retreiveing file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Warnf("error closing file: %s", err)
		}
	}()

	filename := header.Filename
	h.log.Debugf("Uploaded File: %+v, size: %+v, mime: %+v", filename, header.Size, header.Header)

	tempFile, err := os.Create(fmt.Sprintf("%s/%s", workDir, filename))
	if err != nil {
		h.log.Errorf("error creating temporary file for upload, dir: %s, filename: %s", workDir, filename)
		return "", err
	}
	defer func() {
		if err := tempFile.Close(); err != nil {
			h.log.Warnf("error closing file: %s", err)
		}
	}()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	if _, err := tempFile.Write(fileBytes); err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}
	return filename, nil
}

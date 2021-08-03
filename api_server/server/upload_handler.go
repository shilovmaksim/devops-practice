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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

var (
	uiUploadDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ui_upload_duration_time_seconds",
		Help:    "Duration of files uploading",
		Buckets: prometheus.DefBuckets,
	})
	compresionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "compression_time_seconds",
		Help:    "Duration of files compressing",
		Buckets: prometheus.DefBuckets,
	})
	storageUploadDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "storage_upload_duration_time_seconds",
		Help:    "Duration of files uploading files to the storage",
		Buckets: prometheus.DefBuckets,
	})
	optimizationRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "optimization_request_duration_time_seconds",
		Help:    "Duration of optimization request",
		Buckets: prometheus.DefBuckets,
	})
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

// Upload files
// @Summary Upload files and start optimization process
// @Description Upload files for optimization script and request an optimization run
// @ID upload-handler
// @Accept  multipart/form-data
// @Produce  json
// @Param   file formData file true  "filename"
// @Success 200 {object} models.UploadResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/upload [post]
func (h *UploadHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {

	// Create a temporary directory
	env := environment.New(os.TempDir(), envPrefix)
	if err := env.CreateTempDir(); err != nil {
		h.log.Errorf("error creating tempdir: %s", err)
		writeResponse(writer, models.NewErrorResponse(err.Error()), http.StatusInternalServerError, h.log)
	}
	defer func() {
		if err := env.CleanUp(); err != nil {
			h.log.Warnf("error cleaning up temporary directory %s, error: %s", env.Dir(), err)
		}
	}()

	// Upload files from the UI
	h.log.Debugf("Start uploading files")
	timer := prometheus.NewTimer(uiUploadDuration)
	filenames, err := h.serveFileUpload(env.Dir(), r)
	timer.ObserveDuration()
	if err != nil {
		writeResponse(writer, models.NewErrorResponse(ErrMsgInternal), http.StatusBadRequest, h.log)
		return
	}

	// Create zip archive from user provided files
	h.log.Debugf("Add files to zip archive: '%s'", filenames)
	timer = prometheus.NewTimer(compresionDuration)
	absPathToArch := path.Join(env.Dir(), (environment.Filename)(inputFileName).WithUnixSuffix())
	archFilesPath, err := compressor.Compress(context.Background(), absPathToArch, env.Dir(), filenames...)
	timer.ObserveDuration()
	if err != nil {
		writeResponse(writer, models.NewErrorResponse(ErrMsgInternal), http.StatusBadRequest, h.log)
		return
	}

	// Upload files to the bucket
	timer = prometheus.NewTimer(storageUploadDuration)
	filename := filepath.Base(archFilesPath)
	timer.ObserveDuration()
	if _, err := h.storage.UploadFiles(env.Dir(), filename); err != nil {
		writeResponse(writer, models.NewErrorResponse("error uploading archive to the bucket"), http.StatusBadRequest, h.log)
		return
	}

	// POST a request to the optimization service
	h.log.Debugf("Post a request to the optimization service, filename: %s", filename)
	timer = prometheus.NewTimer(optimizationRequestDuration)
	resp, err := h.client.PostOptimize(filename)
	timer.ObserveDuration()
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

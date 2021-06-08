package optimizer

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/python"
	"github.com/cxrdevelop/optimization_engine/pkg/compressor"
	"github.com/cxrdevelop/optimization_engine/pkg/environment"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/cxrdevelop/optimization_engine/pkg/storage"
)

const (
	scriptResultFilename = "def_output.csv"
	uploadPrefix         = "opt_result"
	scriptDirectory      = "script"
)

var _ Optimizer = (*RackOptimizer)(nil)

type RackOptimizer struct {
	wrapper *python.Wrapper
	storage storage.Storage
	workDir string
	prefix  string
	log     *logger.Logger
}

func NewRackOptimizer(wrapper *python.Wrapper, storage storage.Storage, workDir string, prefix string, log *logger.Logger) *RackOptimizer {
	return &RackOptimizer{
		wrapper: wrapper,
		storage: storage,
		workDir: workDir,
		prefix:  prefix,
		log:     log,
	}
}

func (r *RackOptimizer) Execute(filename string) (*Result, error) {
	// Create temp dir
	env := environment.New(r.workDir, r.prefix)
	if err := env.CreateTempDir(); err != nil {
		return nil, ErrEnvCreate
	}
	// Create a dir for script in temp dir
	scriptWorkDir := filepath.Join(env.Dir(), scriptDirectory)
	if err := os.MkdirAll(scriptWorkDir, os.ModePerm); err != nil {
		return nil, ErrEnvCreate
	}
	defer func() {
		if err := env.CleanUp(); err != nil {
			r.log.Errorf("error cleaning up temporary directory %s, error: %s", env.Dir(), err)
		}
	}()

	// Download files
	if err := r.storage.DownloadFiles(env.Dir(), filename); err != nil {
		r.log.Errorf("error downloading files: %s", err)
		return nil, ErrDownload
	}

	// Decompress downloaded zip archive into the workDir
	if err := compressor.Decompress(context.Background(), path.Join(env.Dir(), filename), scriptWorkDir); err != nil {
		r.log.Errorf("error decompressing files: %s", err)
		return nil, ErrDecompress
	}

	// Get filenames
	filenames, err := getFilenames(scriptWorkDir)
	if err != nil {
		r.log.Errorf("error reading script directory: %s", err)
		return nil, ErrInternal
	}

	// Execute script
	scriptRes := r.wrapper.Optimize(scriptWorkDir, filenames...)
	if scriptRes.ExitCode != 0 {
		r.log.Errorf("optimization script error: exit code: %d", scriptRes.ExitCode)
		return nil, ErrOptimize
	}

	// Compress the result
	absPathToArch := path.Join(env.Dir(), (environment.Filename)(uploadPrefix).WithUnixSuffix())
	archPath, err := compressor.Compress(context.Background(), absPathToArch, scriptWorkDir, scriptResultFilename)
	if err != nil {
		r.log.Errorf("error compressing file: '%s', error: %s", scriptResultFilename, err)
		return nil, ErrCompress
	}

	// Upload the compressed file
	uploadRes, err := r.storage.UploadFiles(env.Dir(), filepath.Base(archPath))
	if err != nil {
		r.log.Errorf("error uploading files: %s", err)
		return nil, ErrUpload
	}

	return &Result{
		Filename:      uploadRes[0].Filename,
		Location:      uploadRes[0].Location,
		ETag:          uploadRes[0].ETag,
		ExecutionTime: scriptRes.ExecutionTime,
	}, nil
}

func getFilenames(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	filenames := make([]string, 0, len(files))
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames, nil
}

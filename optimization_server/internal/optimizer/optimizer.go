package optimizer

import (
	"errors"
	"time"
)

var (
	ErrEnvCreate  = errors.New("error creating tempdir")
	ErrDownload   = errors.New("download error")
	ErrDecompress = errors.New("decompression error")
	ErrInternal   = errors.New("internal error")
	ErrCompress   = errors.New("compression error")
	ErrOptimize   = errors.New("optimization script error")
	ErrUpload     = errors.New("upload error")
)

type Result struct {
	Filename      string
	Location      string
	ETag          string
	ExecutionTime time.Duration
}

type Optimizer interface {
	Execute(filename string) (*Result, error)
}

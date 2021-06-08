package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

type FSStorage struct {
	bucket string
	log    *logger.Logger
}

var _ Storage = (*FSStorage)(nil)

// NewFSStorage creates a filesystem storage or panics if provided folder is invalid
func NewFSStorage(bucket string, log *logger.Logger) *FSStorage {
	if fi, err := os.Stat(bucket); err != nil {
		panic(fmt.Errorf("error accessing path: %w", err))
	} else if !fi.IsDir() {
		panic(fmt.Errorf("path is not a directory"))
	}

	return &FSStorage{
		bucket: bucket,
		log:    log,
	}
}

// DownloadFiles function takes S3 bucket name and a remote filename.
// It creates a new file with the same name in the workDir. If download fails the file will be empty.
func (s *FSStorage) DownloadFiles(dir string, paths ...string) error {
	for _, path := range paths {
		if err := s.download(dir, path, path); err != nil {
			return err
		}
	}
	return nil
}

// UploadFiles function takes dir name and a local filename.
// It takes a file from the workDir and uploads it to the bucket with the same name.
func (s *FSStorage) UploadFiles(dir string, paths ...string) ([]UploadResult, error) {
	res := make([]UploadResult, 0, len(paths))
	for _, path := range paths {
		if r, err := s.upload(dir, path, path); err != nil {
			return nil, err
		} else {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (s *FSStorage) download(dir string, remoteFilename string, localFilename string) error {
	if localFilename == "" {
		localFilename = remoteFilename
	}
	src := getFilePath(s.bucket, remoteFilename)
	dst := getFilePath(dir, localFilename)
	s.log.Debugf("Downloading file '%s' to '%s", src, dst)

	if err := s.copyFile(src, dst); err != nil {
		return fmt.Errorf("unable to download '%s' from bucket '%s' with '%w'", remoteFilename, s.bucket, err)
	}

	return nil
}

func (s *FSStorage) upload(dir string, localFilename string, remoteFilename string) (*UploadResult, error) {
	if remoteFilename == "" {
		remoteFilename = localFilename
	}
	path := getFilePath(dir, localFilename)
	s.log.Debugf("Open file '%s'", path)

	file, err := os.Open(path)
	if err != nil {
		s.log.Debugf("failed to open file '%s', error: '%s", path, err)
		return nil, fmt.Errorf("failed to open file '%s', error: '%w", path, err)
	}
	defer s.closeFile(file)

	s.log.Debugf("Uploading file '%s' to the bucket...", remoteFilename)

	tmp := getFilePath(s.bucket, remoteFilename)
	if err := s.copyFile(path, tmp); err != nil {
		return nil, fmt.Errorf("unable to upload local file '%s' to local bucket '%s', remote filename: '%s', error: '%w'", localFilename, s.bucket, remoteFilename, err)
	}
	location, err := filepath.Abs(tmp)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path for filename: '%s', error: '%w'", tmp, err)
	}

	return &UploadResult{
		Filename: remoteFilename,
		Location: location,
		ETag:     "",
	}, nil
}

func (c *FSStorage) closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		c.log.Warnf("error closing file: %s", err)
	}
}

func (s *FSStorage) copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.closeFile(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer s.closeFile(destination)
	_, err = io.Copy(destination, source)
	return err
}

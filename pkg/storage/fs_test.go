package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const bucket = "test_bucket"

func TestFSStorage_Creation(t *testing.T) {
	assert.Panics(t, func() { NewFSStorage("some-non-existing-dir", logger.NewTestLogger()) })
	assert.Panics(t, func() { NewFSStorage("mock_client_test.go", logger.NewTestLogger()) })
	assert.NotPanics(t, func() { NewFSStorage(".", logger.NewTestLogger()) })
}

func TestFSStorage_Upload(t *testing.T) {
	createTestBucket(t)
	defer removeTestBucket(t)

	var mock *FSStorage = nil
	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })
	buffer := []byte("123")

	tmpDir := os.TempDir()
	file, err := os.CreateTemp(tmpDir, "test_upload_")
	assert.NoError(t, err)
	_, err = file.Write(buffer)
	assert.NoError(t, err)
	fi, err := file.Stat()
	assert.NoError(t, err)
	fileName := fi.Name()

	assert.NoError(t, file.Close())
	defer removeFile(t, file.Name())

	res, err := mock.UploadFiles(tmpDir, fileName)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Filename, fileName)
	assert.True(t, filepath.IsAbs(res[0].Location))

	filePath := fmt.Sprintf("%s/%s", bucket, fileName)
	file, err = os.Open(filePath)
	assert.NoError(t, err)

	defer removeFile(t, filePath)

	n, err := file.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, len(buffer), n)
	assert.NoError(t, file.Close())
}

func TestFSStorage_Download(t *testing.T) {
	createTestBucket(t)
	defer removeTestBucket(t)

	// Put a file into the test bucket
	file, err := os.CreateTemp(bucket, "test_download_")
	assert.NoError(t, err)
	buffer := []byte("123")
	_, err = file.Write(buffer)
	assert.NoError(t, err)
	fi, err := file.Stat()
	assert.NoError(t, err)
	fileName := fi.Name()

	filePath := fmt.Sprintf("%s/%s", bucket, fileName)
	assert.NoError(t, file.Close())

	defer removeFile(t, filePath)

	var mock *FSStorage = nil
	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })

	tmpDir := os.TempDir()
	err = mock.DownloadFiles(tmpDir, fileName)
	assert.NoError(t, err)

	filePath = fmt.Sprintf("%s/%s", tmpDir, fileName)
	file, err = os.Open(filePath)
	assert.NoError(t, err)

	defer removeFile(t, filePath)

	n, err := file.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, len(buffer), n)
	assert.NoError(t, file.Close())
}

func TestFSStorage_UploadFails(t *testing.T) {
	// create a folder to serve as a test bucket
	createTestBucket(t)
	defer removeTestBucket(t)

	var mock *FSStorage = nil
	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })
	_, err := mock.UploadFiles("", "1")
	assert.Error(t, err)

	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })
	_, err = mock.UploadFiles("", "")
	assert.Error(t, err)
}

func TestFSStorage_DownloadFails(t *testing.T) {
	// create a folder to serve as a test bucket
	createTestBucket(t)
	defer removeTestBucket(t)

	var mock *FSStorage = nil
	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })
	err := mock.DownloadFiles("", "1")
	assert.Error(t, err)

	assert.NotPanics(t, func() { mock = NewFSStorage(bucket, logger.NewTestLogger()) })
	err = mock.DownloadFiles("", "")
	assert.Error(t, err)
}

func removeFile(t *testing.T, path string) {
	err := os.Remove(path)
	assert.NoError(t, err)
}

func createTestBucket(t *testing.T) {
	err := os.Mkdir(bucket, os.ModePerm)
	assert.NoError(t, err)
}

func removeTestBucket(t *testing.T) {
	err := os.RemoveAll(bucket)
	assert.NoError(t, err)
}

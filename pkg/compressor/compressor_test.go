package compressor

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tmpFolder        = "folder_for_arch"
	decompressFolder = "decompress"
)

var archPath = path.Join(tmpFolder, "arch")

func TestCompressFolder_Success(t *testing.T) {
	assert.NoError(t, os.Mkdir(tmpFolder, os.ModePerm))
	defer func() { assert.NoError(t, os.RemoveAll(tmpFolder)) }()

	name, err := Compress(context.Background(), archPath, "test_files", "file1.csv", "file2.csv")
	assert.NoError(t, err)
	assert.Equal(t, name, path.Join(tmpFolder, "arch.tar.gz"))
}

func TestCompressFolder_Error_EmptyPath(t *testing.T) {
	_, err := Compress(context.Background(), "", "workDir")
	assert.Error(t, err)
	assert.Equal(t, err, ErrEmptyPath)
}

func TestCompressFolder_Error_NoFiles(t *testing.T) {
	_, err := Compress(context.Background(), "path", "workDir")
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoFiles)
}

func TestGzipCompressor_Compress_Error(t *testing.T) {
	_, err := Compress(context.Background(), archPath, "test_files", "non existing path to trigger compression error.csv")
	assert.Error(t, err)
}

func TestDecompressFolder_Success(t *testing.T) {

	assert.NoError(t, os.Mkdir(tmpFolder, os.ModePerm))
	defer func() { assert.NoError(t, os.RemoveAll(tmpFolder)) }()

	name := compress(t)

	decFolder := path.Join(tmpFolder, decompressFolder)
	assert.NoError(t, os.MkdirAll(decFolder, os.ModePerm))

	err := Decompress(context.Background(), name, decFolder)
	assert.NoError(t, err)

}

func TestDecompressFolder_Error_EmptyPath(t *testing.T) {
	err := Decompress(context.Background(), "", "folder")
	assert.Error(t, err)
	assert.Equal(t, err, ErrEmptyPath)
}

func TestDecompressFolder_Error_EmptyFolder(t *testing.T) {
	err := Decompress(context.Background(), "arch.tar.gz", "")
	assert.Error(t, err)
	assert.Equal(t, err, ErrEmptyFolderPath)
}

func TestDecompressFolder_Error_WrongFolder(t *testing.T) {

	assert.NoError(t, os.Mkdir(tmpFolder, os.ModePerm))
	defer func() { assert.NoError(t, os.RemoveAll(tmpFolder)) }()

	name := compress(t)

	err := Decompress(context.Background(), name, "non existing path to trigger compression error")
	assert.Error(t, err)
}

func TestDecompressFolder_Error_WrongArch(t *testing.T) {

	err := Decompress(context.Background(), "no existing arch name", "test_files")
	assert.Error(t, err)
}

func compress(t *testing.T) string {
	name, err := Compress(context.Background(), archPath, "test_files", "file1.csv", "file2.csv")
	assert.NoError(t, err)
	assert.Equal(t, name, path.Join(tmpFolder, "arch.tar.gz"))
	return name
}

package compressor

import (
	"context"
	"fmt"
	"os/exec"
)

var (
	ErrEmptyPath       = fmt.Errorf("path can't be empty")
	ErrEmptyFolderPath = fmt.Errorf("path to folder can't be empty")
	ErrNoFiles         = fmt.Errorf("no files were provided")
)

// Compress function compresses input files using system tar and gzip.
// path sets the archive name without '.tar.gz' suffix.
// workDir sets the working directory.
// filenames set the names to be compressed. Can't be empty.
// Function returns the name of the created archive by adding 'tar.gz' suffix to the provided path and an error message.
func Compress(ctx context.Context, path string, workDir string, filenames ...string) (archPath string, err error) {
	if path == "" {
		return "", ErrEmptyPath
	}
	if len(filenames) == 0 {
		return "", ErrNoFiles
	}
	archPath = fmt.Sprintf("%s.tar.gz", path)

	params := []string{"-C", workDir, "-czf", archPath}
	params = append(params, filenames...)
	cmd := exec.CommandContext(ctx, "tar", params...)
	if _, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("error compressing files: '%w'", err)
	}
	return archPath, nil
}

// Decompress function decompresses the archive into the provided directory using system tar and gzip.
// path sets the archive name.
// folder sets folder in which the archive will be decompressed. Can't be empty.
func Decompress(ctx context.Context, path string, folder string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if folder == "" {
		return ErrEmptyFolderPath
	}

	params := []string{"-xzf", path, "-C", folder}
	cmd := exec.CommandContext(ctx, "tar", params...)
	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("error decompressing files: '%w', output: '%s'", err, out)
	}
	return nil
}

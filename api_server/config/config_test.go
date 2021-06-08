package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_File(t *testing.T) {
	assert.NoError(t, os.Setenv("LOG_PATH", "test_log_path"))
	cfg, err := ReadConfig("test_data/test.yml")
	assert.NoError(t, err)
	assert.Equal(t, cfg.Application.Port, "8085")
	assert.Equal(t, cfg.Application.LogPath, "test_log_path")
	assert.Equal(t, cfg.Application.LogLevel, "debug")
	assert.Equal(t, strings.ToLower(cfg.Storage.Type), "local")
	os.Clearenv()
}

func TestConfig_Environment(t *testing.T) {
	assert.NoError(t, os.Setenv("STORAGE_TYPE", "S3"))
	cfg, err := ReadConfig("test_data/test.yml")
	assert.NoError(t, err)
	assert.Equal(t, strings.ToLower(cfg.Storage.Type), "s3")
	os.Clearenv()
}

func TestConfig_NotFound(t *testing.T) {
	cfg, err := ReadConfig("no_file.yml")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

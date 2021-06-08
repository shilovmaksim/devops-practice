package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

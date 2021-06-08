package storage

import (
	"os"
	"testing"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestS3Storage_Creation(t *testing.T) {
	assert.Panics(t, func() { NewS3Storage("region", "dir", logger.NewTestLogger()) })
	assert.NoError(t, os.Setenv("AWS_ACCESS_KEY_ID", "1"))
	assert.NoError(t, os.Setenv("AWS_SECRET_ACCESS_KEY", "2"))
	assert.NotPanics(t, func() { NewS3Storage("region", "dir", logger.NewTestLogger()) })
}

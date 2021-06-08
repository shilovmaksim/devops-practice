package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/models"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	responseData := models.NewHealthResponse()
	recorder := httptest.NewRecorder()

	writeResponse(recorder, responseData, http.StatusOK, logger.NewTestLogger())
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"health\":true}", recorder.Body.String())
	assert.Equal(t, 200, recorder.Code)
}

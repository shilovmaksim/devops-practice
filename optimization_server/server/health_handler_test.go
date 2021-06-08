package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	health := NewHealthHandler(logger.NewTestLogger())
	handle := func(w http.ResponseWriter, r *http.Request) {
		health.ServeHTTP(w, r)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/api/v1/health", nil)
	assert.NoError(t, err)

	r := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)

	handler.ServeHTTP(r, req)

	assert.Equal(t, r.Code, http.StatusOK)
	assert.Equal(t, r.Body.String(), `{"health":true}`)

}

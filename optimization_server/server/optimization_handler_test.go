package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/optimizer"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestOptimizationHandler(t *testing.T) {

	testCases := []struct {
		name           string
		inputJson      string
		expectedStatus int
		outputJson     string
		optimizer      optimizer.Optimizer
	}{
		{
			name:           "empty",
			inputJson:      "",
			expectedStatus: http.StatusBadRequest,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgEmptyRequest),
			optimizer:      nil,
		},
		{
			name:           "invalid json",
			inputJson:      "123",
			expectedStatus: http.StatusBadRequest,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgJsonParse),
			optimizer:      nil,
		},
		{
			name:           "empty args",
			inputJson:      `{"args":[]}`,
			expectedStatus: http.StatusBadRequest,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErMsgJsonValidate),
			optimizer:      nil,
		},
		{
			name:           "empty arg",
			inputJson:      `{"args":["file1.csv",""]}`,
			expectedStatus: http.StatusBadRequest,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErMsgJsonValidate),
			optimizer:      nil,
		},
		{
			name:           "pseudo error mock, internal error",
			inputJson:      `{"filename":"1"}`,
			expectedStatus: http.StatusInternalServerError,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgInternal),
			optimizer: func() optimizer.Optimizer {
				opt := optimizer.NewMockOptimizer()
				opt.On("Execute", "1").Return(&optimizer.Result{}, optimizer.ErrEnvCreate)
				return opt
			}(),
		},
		{
			name:           "pseudo error mock, download error",
			inputJson:      `{"filename":"1"}`,
			expectedStatus: http.StatusInternalServerError,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgDownload),
			optimizer: func() optimizer.Optimizer {
				opt := optimizer.NewMockOptimizer()
				opt.On("Execute", "1").Return(&optimizer.Result{}, optimizer.ErrDownload)
				return opt
			}(),
		},
		{
			name:           "pseudo error mock, script error",
			inputJson:      `{"filename":"1"}`,
			expectedStatus: http.StatusInternalServerError,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgScript),
			optimizer: func() optimizer.Optimizer {
				opt := optimizer.NewMockOptimizer()
				opt.On("Execute", "1").Return(&optimizer.Result{}, optimizer.ErrOptimize)
				return opt
			}(),
		},
		{
			name:           "pseudo error mock, upload error",
			inputJson:      `{"filename":"1"}`,
			expectedStatus: http.StatusInternalServerError,
			outputJson:     fmt.Sprintf(`{"text":"%s"}`, ErrMsgUpload),
			optimizer: func() optimizer.Optimizer {
				opt := optimizer.NewMockOptimizer()
				opt.On("Execute", "1").Return(&optimizer.Result{}, optimizer.ErrUpload)
				return opt
			}(),
		},
		{
			name:           "success",
			inputJson:      `{"filename":"1"}`,
			expectedStatus: http.StatusOK,
			outputJson:     `{"filename":"1","location":"1","etag":"1","executionTime":0}`,
			optimizer: func() optimizer.Optimizer {
				opt := optimizer.NewMockOptimizer()
				opt.On("Execute", "1").Return(&optimizer.Result{
					Filename:      "1",
					Location:      "1",
					ETag:          "1",
					ExecutionTime: 1,
				}, nil)
				return opt
			}(),
		},
		/*
			{
				name:           "pseudo error mock, no error",
				inputJson:      `{"args":["1", "2", "3", "4", "5"]}`,
				expectedStatus: http.StatusOK,
				outputJson:     `{"filename":"1","location":"1","etag":"1","executionTime":0}`,
			},*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test := func(req *http.Request, err error) {
				optimize := NewOptimizationHandler(tc.optimizer, logger.NewTestLogger())
				handle := func(w http.ResponseWriter, r *http.Request) {
					optimize.ServeHTTP(w, r)
				}

				assert.NoError(t, err)
				r := httptest.NewRecorder()
				handler := http.HandlerFunc(handle)
				handler.ServeHTTP(r, req)
				assert.Equal(t, r.Code, tc.expectedStatus)
				assert.Equal(t, r.Body.String(), tc.outputJson)
			}

			if len(tc.inputJson) == 0 {
				test(http.NewRequestWithContext(context.Background(), "GET", "/api/v1/optimize", nil))
				return
			}

			buf := bytes.NewBuffer([]byte(tc.inputJson))
			test(http.NewRequestWithContext(context.Background(), "GET", "/api/v1/optimize", buf))
		})
	}
}

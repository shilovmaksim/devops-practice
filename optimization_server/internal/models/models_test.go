package models

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	assert.Equal(t, NewHealthResponse().Health, true)
}

func TestError(t *testing.T) {
	assert.Equal(t, NewErrorResponse("errMsg").Text, "errMsg")
}

func TestValidation_OptimizationRequest(t *testing.T) {
	testCases := []struct {
		name     string
		req      OptimizationRequest
		expected url.Values
	}{
		{
			name: "empty name",
			req: OptimizationRequest{
				Filename: "",
			},
			expected: url.Values{"filename": []string{"not set"}},
		},
		{
			name: "valid",
			req: OptimizationRequest{
				Filename: "1.tar.gz",
			},
			expected: url.Values{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.req.Validate()
			assert.Equal(t, len(res), len(tc.expected))

			for k, v := range tc.expected {
				assert.True(t, len(res[k]) == len(v))
				for i := 0; i < len(res[k]); i++ {
					assert.Equal(t, res[k][i], v[i])
				}
			}
		})
	}
}

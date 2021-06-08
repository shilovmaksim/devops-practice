package python

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWrapper(t *testing.T) {
	testCases := []struct {
		name      string
		timeout   time.Duration
		req       []string
		scriptDir string
		result    OptimizationScriptResult
	}{
		{
			name:    "normal run",
			timeout: 5000 * time.Millisecond,
			req:     []string{},
			result: OptimizationScriptResult{
				ExitCode:      0,
				ShellOutput:   "",
				ScriptOutput:  "",
				ExecutionTime: 0,
			},
		},
		{
			name:    "error in script",
			timeout: 5000 * time.Millisecond,
			req:     []string{"error"},
			result: OptimizationScriptResult{
				ExitCode:      14,
				ShellOutput:   "smth",
				ScriptOutput:  "",
				ExecutionTime: 0,
			},
		},
		{
			name:    "successful run",
			timeout: 5000 * time.Millisecond,
			req:     []string{"success"},
			result: OptimizationScriptResult{
				ExitCode:      0,
				ShellOutput:   "",
				ScriptOutput:  "",
				ExecutionTime: 0,
			},
		},
		{
			name:    "print message to std output",
			timeout: 5000 * time.Millisecond,
			req:     []string{"print", "message"},
			result: OptimizationScriptResult{
				ExitCode:      0,
				ShellOutput:   "",
				ScriptOutput:  "message",
				ExecutionTime: 0,
			},
		},
		{
			name:    "test execution time",
			timeout: 5000 * time.Millisecond,
			req:     []string{"sleep", "100"},
			result: OptimizationScriptResult{
				ExitCode:      0,
				ShellOutput:   "",
				ScriptOutput:  "",
				ExecutionTime: 10 * time.Millisecond,
			},
		},
		{
			name:    "test timeout",
			timeout: 5 * time.Millisecond,
			req:     []string{"sleep", "100"},
			result: OptimizationScriptResult{
				ExitCode:      -1,
				ShellOutput:   "smth",
				ScriptOutput:  "",
				ExecutionTime: 1 * time.Millisecond,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := testWrapper(t, tc.req, tc.timeout)
			assert.Equal(t, tc.result.ExitCode, result.ExitCode)
			assert.Equal(t, true, tc.result.ExecutionTime <= result.ExecutionTime)
			assert.Equal(t, true, (len(tc.result.ScriptOutput) > 0) == (len(result.ScriptOutput) > 0), "Script output mismatch")
			assert.Equal(t, true, (len(tc.result.ShellOutput) > 0) == (len(result.ShellOutput) > 0), "Shell output mismatch")
		})
	}
}

func testWrapper(t *testing.T, req []string, timeout time.Duration) *OptimizationScriptResult {

	scriptPath, err := filepath.Abs("mock_scripts/main.py")
	assert.NoError(t, err)

	resp := NewWrapper(scriptPath, timeout, logger.NewTestLogger()).Optimize("", req...)
	return resp
}

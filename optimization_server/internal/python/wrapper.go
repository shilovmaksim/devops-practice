package python

import (
	"context"
	"os/exec"
	"time"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

const commandName = "python3"

type Wrapper struct {
	scriptPath string
	timeout    time.Duration
	log        *logger.Logger
}

type OptimizationScriptResult struct {
	ExitCode      int
	ShellOutput   string
	ScriptOutput  string
	ExecutionTime time.Duration
}

func NewWrapper(scriptPath string, timeout time.Duration, log *logger.Logger) *Wrapper {
	return &Wrapper{
		scriptPath: scriptPath,
		timeout:    timeout,
		log:        log,
	}
}

// Optimize takes input args and runs a script which resides at scriptPath.
// If the path is not absolute, the function will attempt to run the script relatively to workDir
// The workDir is used as a working directory for the script
func (w *Wrapper) Optimize(workDir string, scriptArgs ...string) *OptimizationScriptResult {
	result := OptimizationScriptResult{}

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	// append script name as the first argument
	args := []string{w.scriptPath}
	args = append(args, scriptArgs...)

	cmd := w.command(ctx, workDir, args)

	start := time.Now()
	out, err := cmd.CombinedOutput()

	w.log.Debugf("executed %q %s -> %q", cmd.Path, cmd.Args, out)
	result.ExecutionTime = time.Since(start)
	result.ScriptOutput = string(out)
	result.ExitCode = cmd.ProcessState.ExitCode()

	if err != nil {
		result.ShellOutput = err.Error()
	} else if err = ctx.Err(); err != nil {
		result.ShellOutput = err.Error()
	}

	return &result
}

func (w *Wrapper) command(ctx context.Context, workDir string, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, commandName, args...)
	cmd.Dir = workDir
	return cmd
}

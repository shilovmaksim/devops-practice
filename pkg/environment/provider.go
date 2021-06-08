package environment

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Provider struct {
	workDir  string
	prefix   string
	tempPath string // if tempPath is empty, no folder was created
}

func New(workDir string, prefix string) *Provider {
	return &Provider{
		workDir: workDir,
		prefix:  prefix,
	}
}

func (ep *Provider) CreateTempDir() error {
	dir, err := ioutil.TempDir(ep.workDir, ep.prefix)
	if err != nil {
		return err
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	ep.tempPath = path
	return nil
}

func (ep *Provider) Dir() string {
	return ep.tempPath
}

func (ep *Provider) CleanUp() error {
	if ep.tempPath == "" {
		return nil
	}
	if err := os.RemoveAll(ep.tempPath); err != nil {
		return err
	}
	ep.tempPath = ""
	return nil
}

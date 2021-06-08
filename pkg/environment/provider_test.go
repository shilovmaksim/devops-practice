package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvProvider_UniqueDirs(t *testing.T) {

	env1 := New("./", "test_prefix")
	assert.NoError(t, env1.CreateTempDir())
	dir1 := env1.Dir()
	assert.DirExists(t, dir1)

	env2 := New("./", "test_prefix")
	assert.NoError(t, env2.CreateTempDir())
	dir2 := env2.Dir()
	assert.DirExists(t, dir2)

	assert.NotEqual(t, dir1, dir2)
	assert.NoError(t, env1.CleanUp())
	assert.NoError(t, env2.CleanUp())

}

func TestEnvProvider_Expected(t *testing.T) {
	env := New("./", "test_prefix")
	assert.NoError(t, env.CreateTempDir())
	assert.DirExists(t, env.Dir())
	assert.NoError(t, env.CleanUp())
	assert.Equal(t, env.Dir(), "")
	assert.NoError(t, env.CleanUp())
	assert.NoDirExists(t, env.Dir())
}

func TestEnvProvider_WrongPrefix(t *testing.T) {
	env := New(".", "/")
	assert.Error(t, env.CreateTempDir())
	assert.Equal(t, env.Dir(), "")
	assert.NoDirExists(t, env.Dir())
	assert.NoError(t, env.CleanUp())
}

func TestEnvProvider_WrongWorkDir(t *testing.T) {
	env := New("\\", "")
	assert.Error(t, env.CreateTempDir())
	assert.Equal(t, env.Dir(), "")
	assert.NoDirExists(t, env.Dir())
	assert.NoError(t, env.CleanUp())
}

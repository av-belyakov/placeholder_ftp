package supportingfunctions_test

import (
	"os"
	"path"
	"testing"

	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	var (
		rootDir   string = "placeholder_ftp"
		createDir string = "tmp_dir"
	)

	err := supportingfunctions.CreateDirectory(rootDir, createDir)
	assert.NoError(t, err)

	newCreatePath, err := supportingfunctions.GetRootPath(rootDir)
	assert.NoError(t, err)

	newwPath := path.Join(newCreatePath, createDir)
	_, err = os.ReadDir(newwPath)
	assert.NoError(t, err)

	err = os.RemoveAll(newwPath)
	assert.NoError(t, err)
}

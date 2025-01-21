package writingfilelimit_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
)

func TestWritingFileLimit(t *testing.T) {
	var (
		pathName string = "../../test/test_files/"
		fileName string = "___test_pcap_file.pcap.txt"
	)

	num, err := supportingfunctions.WritingFileLimit(pathName, fileName, ".limit", 40000)
	assert.NoError(t, err)

	t.Log("writed to file:", num, " byte")

	fileInfo, err := os.Stat(path.Join(pathName, fileName))
	assert.NoError(t, err)

	t.Log("file size:", fileInfo.Size())

	assert.NotEqual(t, fileInfo.Size(), 0)
}

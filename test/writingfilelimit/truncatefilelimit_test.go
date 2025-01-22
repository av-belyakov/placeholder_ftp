package writingfilelimit_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateFileLimit(t *testing.T) {
	//test_pcap_file.pcap.txt

	err := os.Truncate("../test_files/test_pcap_file.pcap.txt", 30000)
	assert.NoError(t, err)
}

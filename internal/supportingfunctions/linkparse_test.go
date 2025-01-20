package supportingfunctions_test

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/stretchr/testify/assert"
)

func TestLinkToFile(t *testing.T) {
	links := []string{
		"ftp://ftp.cloud.gcm/traffic/8030101/1737292823_2025_01_19____16_20_23_739436.pcap",
		"ftp://ftp.cloud.gcm/traffic/8030059/1737270561_2025_01_19____10_09_21_577797.pcap",
		"ftp://ftp.cloud.gcm/traffic/8030030/1737251995_2025_01_19____04_59_55_342570.pcap",
		"ftp://10.8.0.7//traff//430032/1737349487_2025_01_20____08_04_47_561660.pcap",
		"ftp://10.67.1.21/traffic/690026/1737316356_2025_01_19____22_52_36_920916.pcap",
		"ftp://10.67.1.21/traffic/690026/1735695602_2025_01_01____04_40_02_337340.pcap",
	}

	for k, v := range links {
		ok := strings.HasPrefix(v, "ftp://")
		assert.True(t, ok)

		ok = strings.HasSuffix(v, ".pcap")
		assert.True(t, ok)

		result, err := supportingfunctions.LinkParse(v)
		assert.NoError(t, err)

		u := &url.URL{
			Scheme: result.Scheme,
			Host:   result.Host,
			Path:   path.Join(result.Path, result.FileName),
		}

		fmt.Printf("%d.\nold_url:%s\nnew_url:%s\npath:%s\nfile:%s\n", k, v, u, result.Path, result.FileName)

		if k != 3 {
			assert.Equal(t, v, u.String())
		}
	}
}

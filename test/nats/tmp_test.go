package nats_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TmpType struct {
	TaskId     string        `json:"task_id"`
	Source     string        `json:"source"`
	Service    string        `json:"service"`
	Command    string        `json:"command"`
	Parameters TmpParameters `json:"parameters"`
}

type TmpParameters struct {
	Links []string `json:"link"`
}

func TestTmp(t *testing.T) {
	raw := []byte(fmt.Sprintf(`{
	"task_id": "6ffab1ea-27ad-4129-925c-e2680c267d62",
	"source": "rcmspb",
	"service": "placeholder_ftp_client_test",
	"command": "convert_and_copy_file",
	"parameters": {
		"links": [
		  "ftp://zsiem-ftp.rcm.spbfsb.ru/traf/500044/1739437517_2025_02_13____12_05_17_534997.pcap",
		  "ftp://zsiem-ftp.rcm.spbfsb.ru/traf/500036/1739442290_2025_02_13____13_24_50_768481.pcap"
		]
	}
}`,
	))

	tt := TmpType{}
	err := json.Unmarshal(raw, &tt)

	assert.NoError(t, err)
}

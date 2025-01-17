package nats_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestSendMsgToModuleFTP(t *testing.T) {
	var Host string = "nats.cloud.gcm"
	var Port int = 4222

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", Host, Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println(err)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println(err)
	})

	replyTo := nats.NewInbox()
	err = nc.PublishRequest(
		"phftp.commands",
		replyTo,
		[]byte(fmt.Sprintf(`{
			"task_id": "6ffab1ea-27ad-4129-925c-e2680c267d62",
			"source": "gcm",
			"service": "placeholder_ftp_client",
			"command": "convert_and_copy_file",
			"parameters": {
				"path_local_ftp": "/traffic/8030164",
				"path_main_ftp": "/traffic/8030164",
				"files": [
				  "1663128065_2022_09_14____07_01_05_749644.pcap",
				  "1663143227_2022_09_14____11_13_47_575934.pcap" 
				]
			}
		}`)))
	/*[]byte(fmt.Sprintf(`{
		"task_id": "%s",
		"source": "gcm",
		"service": "test_service",
		"command": "convert_and_copy_file",
		"parameters": {
			"path_local_ftp": "/net_traff",
			"path_main_ftp": "/net_traff_txt",
			"files": ["test_pcap_file.pcap", "test_pcap_file_http.pcap"]
		}
	}`, uuid.New().String())))*/
	assert.NoError(t, err)

	sub, err := nc.SubscribeSync(replyTo)
	assert.NoError(t, err)

	msg, err := sub.NextMsg(15 * time.Second)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(msg.Data, &response)
	assert.NoError(t, err)

	fmt.Println("--- RESPOND MSG:")
	for k, v := range response {
		fmt.Printf("%s: %v\n", k, v)
	}

	assert.True(t, true)
}

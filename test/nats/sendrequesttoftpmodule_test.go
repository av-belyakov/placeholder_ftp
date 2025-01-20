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
		"phftp.commands.test",
		replyTo,
		//с моего локального ftp на мой локальный
		/*[]byte(fmt.Sprintf(`{
			"task_id": "6ffab1ea-27ad-4129-925c-e2680c267d62",
			"source": "rcmmsk",
			"service": "placeholder_ftp_client",
			"command": "convert_and_copy_file",
			"parameters": {
				"links": [
				  "ftp://127.0.0.1/folder_one/test_pcap_file.pcap",
				  "ftp://127.0.0.1/folder_one/test_pcap_file1_http.pcap"
				]
			}
		}`*/

		//с ftp-users.cloud.gcm на ftp.cloud.gcm
		[]byte(fmt.Sprintf(`{
			"task_id": "6ffab1ea-27ad-4129-925c-e2680c267d62",
			"source": "gcm",
			"service": "placeholder_ftp_client",
			"command": "convert_and_copy_file",
			"parameters": {
				"links": [
				  "ftp://ftp-users.cloud.gcm/net_fraff/test_pcap_file.pcap",
				  "ftp://ftp-users.cloud.gcm/net_traff/test_pcap_file_http.pcap", 
				  "ftp://ftp-users.cloud.gcm/net_fraff/1616152425_2021_03_19____14_13_45_24636.tdp",
				  "ftp://ftp-users.cloud.gcm/net_traff/1616152317_2021_03_19____14_11_57_59.tdp",
				  "ftp://ftp-users.cloud.gcm/net_traff/1636150859_2021_11_06____01_20_59_187344.pcap"
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

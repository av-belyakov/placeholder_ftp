package nats_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
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
	"task_id": "%s",
	"service": "test_service",
	"command": "copy_file",
	"parameters": {
		"path_local_ftp": "/ftp/someuser/folder_one",
		"path_main_ftp": "/ftp/someuser/folder_two",
		"files": ["book.pdf"]
	}
}`, uuid.New().String())))
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

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	HOST string = "nats.cloud.gcm"
	PORT int    = 4222
)

var (
	chMsg chan bool
	nc    *nats.Conn

	err error
)

func TestMain(m *testing.M) {
	chMsg = make(chan bool)

	nc, err = nats.Connect(
		fmt.Sprintf("%s:%d", HOST, PORT),
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

	go func() {
		nc.Subscribe("phftp.commands", func(msg *nats.Msg) {
			b, err := json.MarshalIndent(msg.Data, "", " ")
			if err != nil {
				fmt.Println("ERRORRRRR:", err)
			}
			fmt.Println("func 'subscriptionHandler', Incoming request:", string(b))

			chMsg <- true
		})
	}()

	os.Exit(m.Run())
}

func TestSendMsgToNats(t *testing.T) {
	t.Run("Тест 1. Отправка данных в NATS", func(t *testing.T) {
		err := nc.Publish("phftp.commands",
			[]byte(fmt.Sprintf(`{
			"task_id": "%s",
			"service": "test_service",
			"command": "copy_file",
			"parameters": {
				"path_local_ftp": "/someuser/folder_one",
				"path_main_ftp": "/someuser/folder_two",
				"files": ["book.pdf"]
			}
		}`, uuid.New().String())))

		assert.NoError(t, err)
	})

	t.Run("Тест 2. Проверка приема сообщения", func(t *testing.T) {
		assert.True(t, <-chMsg)
	})
}

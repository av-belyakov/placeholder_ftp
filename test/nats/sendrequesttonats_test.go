package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/cmd/messagebrokers/natsapi"
	"github.com/av-belyakov/placeholder_ftp/internal/logginghandler"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	Host string = "nats.cloud.gcm"
	Port int    = 4222
)

var (
	chMsg chan bool
	nc    *nats.Conn

	chNatsAPIReq <-chan commoninterfaces.ChannelRequester

	ctx       context.Context
	ctxCancel context.CancelFunc

	err error
)

func TestMain(m *testing.M) {
	chMsg = make(chan bool)
	ctx, ctxCancel = context.WithCancel(context.Background())

	nc, err = nats.Connect(
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

	/*go func() {
		nc.Subscribe("phftp.commands", func(msg *nats.Msg) {
			b, err := json.MarshalIndent(msg.Data, "", " ")
			if err != nil {
				fmt.Println("ERRORRRRR:", err)
			}
			fmt.Println("func 'subscriptionHandler', Incoming request:", string(b))

			chMsg <- true
		})
	}()*/

	logging := logginghandler.New()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-logging.GetChan():
				fmt.Printf("LOG TYPE:'%s', MSG:'%s'\n", msg.GetType(), msg.GetMessage())
			}
		}
	}()

	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(Host),
		natsapi.WithPort(Port),
		natsapi.WithCacheTTL(60),
		natsapi.WithSubListenerCommand("phftp.commands")}
	apiNats, err := natsapi.New(logging, natsOptsAPI...)
	if err != nil {
		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chNatsAPIReq, err = apiNats.Start(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

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
		msg := <-chNatsAPIReq

		data := map[string]interface{}{}
		err := json.Unmarshal(msg.GetData(), &data)
		assert.NoError(t, err)

		fmt.Println("RECEIVED MESSAGE FROM NATS:", data)

		assert.True(t, true)
	})
}

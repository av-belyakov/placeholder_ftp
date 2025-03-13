// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger ci.Logger, nameRegionalObject string, opts ...NatsApiOptions) (*apiNatsModule, error) {
	api := &apiNatsModule{
		cachettl:           10,
		logger:             logger,
		nameRegionalObject: nameRegionalObject,
		sendingChannel:     make(chan ci.ChannelRequester),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// Start инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiNatsModule) Start(ctx context.Context) (<-chan ci.ChannelRequester, error) {
	if ctx.Err() != nil {
		return api.sendingChannel, ctx.Err()
	}

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.host, api.port),
		//имя клиента
		nats.Name("thehivehook"),
		//неограниченное количество попыток переподключения
		nats.MaxReconnects(-1),
		//время ожидания до следующей попытки переподключения (по умолчанию 2 сек.)
		nats.ReconnectWait(3*time.Second),
		//обработка разрыва соединения с NATS
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected (%w)", err)).Error())
		}),
		//обработка переподключения к NATS
		nats.ReconnectHandler(func(c *nats.Conn) {
			api.logger.Send("info", "the connection to NATS has been re-established")
		}))
	if err != nil {
		return api.sendingChannel, supportingfunctions.CustomError(err)
	}

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()
		nc.Close()
	}(ctx, nc)

	api.natsConnection = nc

	//обработчик подписки
	go api.subscriptionHandler(ctx)

	return api.sendingChannel, nil
}

// subscriptionHandler обработчик команд
func (api *apiNatsModule) subscriptionHandler(ctx context.Context) {
	api.natsConnection.Subscribe(api.subscriptions.listenerCommand, func(m *nats.Msg) {
		go api.handlerIncomingJSON(ctx, m)
	})
}

// handlerIncomingJSON обработчик JSON объектов
func (api *apiNatsModule) handlerIncomingJSON(ctx context.Context, m *nats.Msg) {
	rc := RequestCommand{}
	if err := json.Unmarshal(m.Data, &rc); err != nil {
		api.logger.Send("error", supportingfunctions.CustomError(err).Error())

		//сообщение с ошибкой из-за полученного некорректного JSON объекта
		res, err := json.Marshal(MainResponse{Error: "invalid json object received"})
		if err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		if err := m.Respond(res); err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		return
	}

	var response MainResponse
	if rc.TaskId == "" {
		response.Error = "invalid json object received, issue 'task_id' is missing"
	}

	if rc.Command == "" {
		response.Error = fmt.Sprintf("invalid json object received, issue 'command' is missing (task id:'%s')", rc.TaskId)
	}

	if rc.Source == "" {
		response.Error = fmt.Sprintf("invalid json object received, issue 'source' is missing (task id:'%s')", rc.TaskId)
	}

	if response.Error != "" {
		res, err := json.Marshal(response)
		if err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		if err := m.Respond(res); err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		api.logger.Send("error", response.Error)

		return
	}

	//убеждаемся, что входящий запрос действительно предназначен для обработки
	//текущим региональным объектом
	if rc.Source != api.nameRegionalObject {
		api.logger.Send("warning", fmt.Sprintf("source name '%s' does not match the current regional name of the object (task id: '%s')", rc.Source, rc.TaskId))

		return
	}

	//отправляем сообщение информирующее получателя о том что, запрос был принят
	//тем региональным объектом для которого он был предназначени и находится в
	//процессе обработки
	m.Respond([]byte(fmt.Sprintf(`{
		"cm_name": "%s",
		"is_processing": "true"
	  }`, api.nameRegionalObject)))

	go api.handlerIncomingCommands(ctx, rc, m)
}

// handlerIncomingCommands обработчик команд
func (api *apiNatsModule) handlerIncomingCommands(ctx context.Context, rc RequestCommand, m *nats.Msg) {
	chRes := make(chan ci.ChannelResponser)

	ttlTime := (time.Duration(api.cachettl) * time.Second)
	ctxTimeout, ctxTimeoutCancel := context.WithTimeout(ctx, ttlTime)
	defer func(cancel context.CancelFunc, ch chan ci.ChannelResponser) {
		cancel()
		close(ch)
		ch = nil
	}(ctxTimeoutCancel, chRes)

	api.sendingChannel <- &RequestFromNats{
		RequestId:  rc.TaskId,
		Command:    "send_command",
		Order:      rc.Command,
		Data:       m.Data,
		ChanOutput: chRes,
	}

	for {
		select {
		case <-ctxTimeout.Done():
			return

		case msg := <-chRes:
			listProcessedFile := []ProcessedInformation(nil)
			for _, v := range msg.GetData() {
				processedFile := ProcessedInformation{
					Error:               "'no error'",
					LinkOld:             v.GetLinkOld(),
					LinkNew:             v.GetLinkNew(),
					SizeBeforProcessing: v.GetSizeBeforProcessing(),
					SizeAfterProcessing: v.GetSizeAfterProcessing(),
				}

				//ошибки по обработки каждой ссылке
				if v.GetError() != nil {
					processedFile.Error = v.GetError().Error()
				}

				listProcessedFile = append(listProcessedFile, processedFile)
			}

			mainResponse := MainResponse{
				Error:             "'no error'",
				Source:            rc.Source,
				RequestId:         rc.TaskId,
				ListProcessedFile: listProcessedFile,
			}

			//обработка глобальной ошибки
			if msg.GetError() != nil {
				mainResponse.Error = msg.GetError().Error()
			}

			res, err := json.Marshal(mainResponse)
			if err != nil {
				err = fmt.Errorf("%w (task id: '%s')", err, rc.TaskId)
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}

			if err := m.Respond(res); err != nil {
				err = fmt.Errorf("%w (task id: '%s')", err, rc.TaskId)
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}

			api.logger.Send("info", fmt.Sprintf("task id '%s' the command '%s' from service '%s' returned %+v", msg.GetRequestId(), rc.Command, rc.Service, listProcessedFile))
		}
	}
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithCacheTTL устанавливает время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL(v int) NatsApiOptions {
	return func(th *apiNatsModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.cachettl = v

		return nil
	}
}

// WithSubSenderCase устанавливает канал в который будут отправлятся объекты типа 'case'
func WithSubSenderCase(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.senderCase = v

		return nil
	}
}

// WithSubSenderAlert устанавливает канал в который будут отправлятся объекты типа 'alert'
func WithSubSenderAlert(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'sender_alert' cannot be empty")
		}

		n.subscriptions.senderAlert = v

		return nil
	}
}

// WithSubListenerCommand устанавливает канал через которые будут приходить команды для
// выполнения определенных действий в TheHive
func WithSubListenerCommand(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'listener_command' cannot be empty")
		}

		n.subscriptions.listenerCommand = v

		return nil
	}
}

func (mnats *ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan {
	return mnats.chanOutputNATS
}

func (mnats *ModuleNATS) SendingData(data SettingsOutputChan) {
	mnats.chanOutputNATS <- data
}

// WithSubscribers метод добавляет абонентов NATS
//func WithSubscribers(event string, responders []string) NatsApiOptions {
//	return func(n *apiNatsModule) error {
//		if event == "" {
//			return errors.New("the subscriber element 'event' must not be empty")
//		}
//
//		if len(responders) == 0 {
//			return errors.New("the subscriber element 'responders' must not be empty")
//		}
//
//		n.subscribers = append(n.subscribers, SubscriberNATS{
//			Event:      event,
//			Responders: responders,
//		})
//
//		return nil
//	}
//}

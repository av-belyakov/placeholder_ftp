package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule[T any] struct {
	subscriptions subscription
	logger        commoninterfaces.Logger
	//передача запросов из NATS
	sendingChannel chan commoninterfaces.ChannelRequester[T]
	//соединение с NATS
	natsConnection *nats.Conn
	host           string
	port           int
	cachettl       int
}

type subscription struct {
	senderCase      string
	senderAlert     string
	listenerCommand string
}

// NatsApiOptions функциональные опции
type NatsApiOptions[T any] func(*apiNatsModule[T]) error

// ModuleNATS инициализированный модуль
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan //канал для отправки полученных данных из модуля
}

// SettingsOutputChan канал вывода данных из модуля
type SettingsOutputChan struct {
	MsgId       string //id сообщения
	SubjectType string //тип подписки
	Data        []byte //набор данных
}

// SettingsInputChan канал для приема данных в модуль
type SettingsInputChan struct {
	Command, EventId, TaskId string
}

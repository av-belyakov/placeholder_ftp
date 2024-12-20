package natsapi

import (
	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// RequestFromNats структура запроса из модуля
type RequestFromNats[T any] struct {
	RequestId   string                      //id запроса
	ElementType string                      //тип элемента
	RootId      string                      //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId      string                      //идентификатор кейса в TheHive
	Command     string                      //команда
	Order       string                      //распоряжение
	Data        []byte                      //набор данных
	ChanOutput  chan ci.ChannelResponser[T] //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToNats структура ответа в модуля
type ResponsToNats struct {
	StatusCode int    //статус кода ответа
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Errror     error  //описание ошибки
	Data       []byte //набор данных
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	TaskId  string `json:"task_id"` //id задачи
	Service string `json:"service"` //наименование сервиса
	Command string `json:"command"` //команда
}

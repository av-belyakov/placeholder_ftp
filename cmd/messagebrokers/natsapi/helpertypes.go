package natsapi

import (
	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// RequestFromNats структура запроса из модуля
type RequestFromNats struct {
	Data        []byte                   //набор данных
	RequestId   string                   //id запроса
	ElementType string                   //тип элемента
	RootId      string                   //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId      string                   //идентификатор кейса в TheHive
	Command     string                   //команда
	Order       string                   //распоряжение
	ChanOutput  chan ci.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToNats структура ответа в модуля
type ResponseToNats struct {
	Data       []commoninterfaces.FileInformationTransfer //набор данных
	Error      error                                      //описание ошибки
	RequestId  string                                     //уникальный идентификатор ответа (соответствует идентификатору запроса)
	StatusCode int                                        //статус кода ответа
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	TaskId  string `json:"task_id"` //уникальный идентификатор задачи
	Source  string `json:"source"`  //наименование регионального объекта к которому был адресован запрос
	Service string `json:"service"` //наименование сервиса
	Command string `json:"command"` //команда
}

// MainResponse основной ответ, на запрос стороннего сервиса
type MainResponse struct {
	ListProcessedFile []ProcessedFile `json:"list_processed_file"`
	RequestId         string          `json:"request_id"`
	Source            string          `json:"source"`
	Error             string          `json:"error"`
}

// ProcessedFile подробное описание результата по обработке файла
type ProcessedFile struct {
	FileNameOld         string `json:"file_name_old"`
	FileNameNew         string `json:"file_name_new"`
	Error               string `json:"error"`
	SizeBeforProcessing int    `json:"size_befor_processing"`
	SizeAfterProcessing int    `json:"size_after_processing"`
}

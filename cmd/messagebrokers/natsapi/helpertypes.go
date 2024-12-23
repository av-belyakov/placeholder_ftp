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
	RequestId  string                                     //UUID идентификатор ответа (соответствует идентификатору запроса)
	StatusCode int                                        //статус кода ответа
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	TaskId  string `json:"task_id"` //id задачи
	Service string `json:"service"` //наименование сервиса
	Command string `json:"command"` //команда
}

// MainResponse основной ответ, на запрос стороннего сервиса
type MainResponse struct {
	ListProcessedFile []ProcessedFile `json:"list_processed_file"`
	Error             string          `json:"error"`
	RequestId         string          `json:"request_id"`
}

// ProcessedFile подробное описание результата по обработке файла
type ProcessedFile struct {
	FileName            string `json:"file_name"`
	Error               string `json:"error"`
	SizeBeforProcessing int    `json:"size_befor_processing"`
	SizeAfterProcessing int    `json:"size_after_processing"`
}

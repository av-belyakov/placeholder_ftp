package handlers

import "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"

// RequestCopyFileFromFtpServer структура запроса для обработки файлов на FTP сервере
type RequestCopyFileFromFtpServer struct {
	TaskId     string                         `json:"task_id"`    //идентификатор задачи
	Service    string                         `json:"service"`    //наименование сервиса
	Command    string                         `json:"command"`    //наименование команды
	Parameters ParameterCopyFileFromFtpServer `json:"parameters"` //дополнительные параметры
}

type ResultRequestCopyFileFromFtpServer struct {
	Data   []commoninterfaces.FileInformationTransfer `json:"data"`    //содержит данные
	Error  error                                      `json:"error"`   //содержит глобальные ошибки, такие как например, ошибка подключения к ftp серверу
	TaskId string                                     `json:"task_id"` //идентификатор задачи
}

type ProcessedFiles struct {
	Error               error  `json:"error"`                 //ошибка возникшая при обработки файла
	FileName            string `json:"file_name"`             //имя файла
	SizeBeforProcessing int    `json:"size_befor_processing"` //размер файла до обработки
	SizeAfterProcessing int    `json:"size_after_processing"` //размер файла после обработки
}

// ParameterCopyFileFromFtpServer подробные параметры
type ParameterCopyFileFromFtpServer struct {
	PathLocalFtp string   `json:"path_local_ftp"`
	PathMainFtp  string   `json:"path_main_ftp"`
	Files        []string `json:"files"`
}

type FtpHandlerOptions struct {
	TmpDir       string
	ConfLocalFtp commoninterfaces.SimpleNetworkConsumer
	ConfMainFtp  commoninterfaces.SimpleNetworkConsumer
	Logger       commoninterfaces.Logger
}

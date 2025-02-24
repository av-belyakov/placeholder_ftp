package handlers

import "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"

// RequestCopyFileFromFtpServer структура запроса для обработки файлов на FTP сервере
type RequestCopyFileFromFtpServer struct {
	TaskId     string                         `json:"task_id"`    //идентификатор задачи
	Service    string                         `json:"service"`    //наименование сервиса
	Command    string                         `json:"command"`    //наименование команды
	Parameters ParameterCopyFileFromFtpServer `json:"parameters"` //дополнительные параметры
}

// ParameterCopyFileFromFtpServer подробные параметры
type ParameterCopyFileFromFtpServer struct {
	Links []string `json:"links"`
}

type ResultRequestCopyFileFromFtpServer struct {
	Data   []commoninterfaces.LinkInformationTransfer `json:"data"`    //содержит данные
	TaskId string                                     `json:"task_id"` //идентификатор задачи
	Error  error                                      `json:"error"`   //содержит глобальные ошибки, такие как например, ошибка подключения к ftp серверу
}

type FtpHandlerOptions struct {
	ConfLocalFtp         commoninterfaces.SimpleNetworkConsumer
	ConfMainFtp          commoninterfaces.SimpleNetworkConsumer
	Logger               commoninterfaces.Logger
	PathResultDirMainFTP string
	TmpDir               string
	MaxWritingFileLimit  int
}

type ProcessedLink struct {
	linkOld             string //старое имя файла
	linkNew             string //новое имя файла (которое формируется на основе старого, после обработки файла декодером)
	sizeBeforProcessing int    //размер файла до обработки
	sizeAfterProcessing int    //размер файла после обработки
	err                 error  //ошибка возникшая при обработки файла
}

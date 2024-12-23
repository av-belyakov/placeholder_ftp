package handlers

import (
	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// ******* ResultRequestCopyFileFromFtpServer *******

// NewResultRequestCopyFileFromFtpServer результат запроса на копирование файла с ftp сервера
func NewResultRequestCopyFileFromFtpServer() *ResultRequestCopyFileFromFtpServer {
	return &ResultRequestCopyFileFromFtpServer{}
}

// GetRequestId возвращает идентификатор запроса
func (obj *ResultRequestCopyFileFromFtpServer) GetRequestId() string {
	return obj.TaskId
}

// SetRequestId устанавливает идентификатор запроса
func (obj *ResultRequestCopyFileFromFtpServer) SetRequestId(v string) {
	obj.TaskId = v
}

// GetError возвращает глобальную ошибку которая может возникнут при выполнении задачи
func (obj *ResultRequestCopyFileFromFtpServer) GetError() error {
	return obj.Error
}

// SetError устанавливает глобальную ошибку которая может возникнут при выполнении задачи
func (obj *ResultRequestCopyFileFromFtpServer) SetError(v error) {
	obj.Error = v
}

// GetData возвращает данные
func (obj *ResultRequestCopyFileFromFtpServer) GetData() []commoninterfaces.FileInformationTransfer {
	return obj.Data
}

// SetData устанавливает данные
func (obj *ResultRequestCopyFileFromFtpServer) SetData(v []commoninterfaces.FileInformationTransfer) {
	obj.Data = v
}

// ******* ProcessedFiles *******

// NewProcessedFiles описание обработанного файла
func NewProcessedFiles() *ProcessedFiles {
	return &ProcessedFiles{}
}

// GetError возвращает ошибку возникшую при обработки файла
func (obj *ProcessedFiles) GetError() error {
	return obj.Error
}

// SetError устанавливает ошибку возникшую при обработки файла
func (obj *ProcessedFiles) SetError(v error) {
	obj.Error = v
}

// GetFileName возвращает имя файла
func (obj *ProcessedFiles) GetFileName() string {
	return obj.FileName
}

// SetFileName устанавливает имя файла
func (obj *ProcessedFiles) SetFileName(v string) {
	obj.FileName = v
}

// GetSizeBeforProcessing возвращает размер файла перед обработкой
func (obj *ProcessedFiles) GetSizeBeforProcessing() int {
	return obj.SizeBeforProcessing
}

// SetSizeBeforProcessing устанавливает размер файла перед обработкой
func (obj *ProcessedFiles) SetSizeBeforProcessing(v int) {
	obj.SizeBeforProcessing = v
}

// GetSizeAfterProcessing возвращает размер файла после обработки
func (obj *ProcessedFiles) GetSizeAfterProcessing() int {
	return obj.SizeAfterProcessing
}

// SetSizeAfterProcessing устанавливает размер файла после обработки
func (obj *ProcessedFiles) SetSizeAfterProcessing(v int) {
	obj.SizeAfterProcessing = v
}

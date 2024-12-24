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
	return obj.err
}

// SetError устанавливает ошибку возникшую при обработки файла
func (obj *ProcessedFiles) SetError(v error) {
	obj.err = v
}

// GetFileNameOld возвращает старое имя файла
func (obj *ProcessedFiles) GetFileNameOld() string {
	return obj.fileNameOld
}

// SetFileNameOld устанавливает старое имя файла
func (obj *ProcessedFiles) SetFileNameOld(v string) {
	obj.fileNameOld = v
}

// GetFileNameNew возвращает новое имя файла
func (obj *ProcessedFiles) GetFileNameNew() string {
	return obj.fileNameNew
}

// SetFileNameNew устанавливает новое имя файла
func (obj *ProcessedFiles) SetFileNameNew(v string) {
	obj.fileNameNew = v
}

// GetSizeBeforProcessing возвращает размер файла перед обработкой
func (obj *ProcessedFiles) GetSizeBeforProcessing() int {
	return obj.sizeBeforProcessing
}

// SetSizeBeforProcessing устанавливает размер файла перед обработкой
func (obj *ProcessedFiles) SetSizeBeforProcessing(v int) {
	obj.sizeBeforProcessing = v
}

// GetSizeAfterProcessing возвращает размер файла после обработки
func (obj *ProcessedFiles) GetSizeAfterProcessing() int {
	return obj.sizeAfterProcessing
}

// SetSizeAfterProcessing устанавливает размер файла после обработки
func (obj *ProcessedFiles) SetSizeAfterProcessing(v int) {
	obj.sizeAfterProcessing = v
}

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
func (obj *ResultRequestCopyFileFromFtpServer) GetData() []commoninterfaces.LinkInformationTransfer {
	return obj.Data
}

// SetData устанавливает данные
func (obj *ResultRequestCopyFileFromFtpServer) SetData(v []commoninterfaces.LinkInformationTransfer) {
	obj.Data = v
}

// ******* ProcessedLink *******

// NewProcessedLink описание обработанного файла
func NewProcessedLink() *ProcessedLink {
	return &ProcessedLink{}
}

// GetError возвращает ошибку возникшую при обработки файла
func (obj *ProcessedLink) GetError() error {
	return obj.err
}

// SetError устанавливает ошибку возникшую при обработки файла
func (obj *ProcessedLink) SetError(v error) {
	obj.err = v
}

// GetLinkOld возвращает старое имя файла
func (obj *ProcessedLink) GetLinkOld() string {
	return obj.linkOld
}

// SetLinkOld устанавливает старое имя файла
func (obj *ProcessedLink) SetLinkOld(v string) {
	obj.linkOld = v
}

// GetLinkNew возвращает новое имя файла
func (obj *ProcessedLink) GetLinkNew() string {
	return obj.linkNew
}

// SetLinkNew устанавливает новое имя файла
func (obj *ProcessedLink) SetLinkNew(v string) {
	obj.linkNew = v
}

// GetSizeBeforProcessing возвращает размер файла перед обработкой
func (obj *ProcessedLink) GetSizeBeforProcessing() int {
	return obj.sizeBeforProcessing
}

// SetSizeBeforProcessing устанавливает размер файла перед обработкой
func (obj *ProcessedLink) SetSizeBeforProcessing(v int) {
	obj.sizeBeforProcessing = v
}

// GetSizeAfterProcessing возвращает размер файла после обработки
func (obj *ProcessedLink) GetSizeAfterProcessing() int {
	return obj.sizeAfterProcessing
}

// SetSizeAfterProcessing устанавливает размер файла после обработки
func (obj *ProcessedLink) SetSizeAfterProcessing(v int) {
	obj.sizeAfterProcessing = v
}

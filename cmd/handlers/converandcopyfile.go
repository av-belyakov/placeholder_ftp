package handlers

import (
	"context"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// HandlerConvertAndCopyFile обработчик копирования и преобразования ффайлов
func (opts FtpHandlerOptions) HandlerConvertAndCopyFile(ctx context.Context, req commoninterfaces.ChannelRequester) {

	//здесь все тоже самое что и в обработчике HandlerCopyFile
	// однако здесь еще необходимо выполнить обработку получаемых
	// pcap файлов
}

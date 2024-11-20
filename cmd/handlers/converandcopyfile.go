package handlers

import (
	"context"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

func HandlerConvertAndCopyFile(
	ctx context.Context,
	req commoninterfaces.ChannelRequester,
	confLocalFTP commoninterfaces.SimpleNetworkConsumer,
	confMainFTP commoninterfaces.SimpleNetworkConsumer,
	logger commoninterfaces.Logger) {

	//здесь все тоже самое что и в обработчике HandlerCopyFile
	// однако здесь еще необходимо выполнить обработку получаемых
	// pcap файлов
}

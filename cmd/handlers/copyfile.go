package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

// /HandlerCopyFile обработчик копирования файлов
func (opts FtpHandlerOptions) HandlerCopyFile(ctx context.Context, req commoninterfaces.ChannelRequester) {
	// исходный ftp сервер
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfLocalFtp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

		return
	}

	request := RequestCopyFileFromFtpServer{}
	if err := json.Unmarshal(req.GetData(), &request); err != nil {
		_, f, l, _ := runtime.Caller(0)
		opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

		return
	}

	for _, file := range request.Parameters.Files {
		b, _, err := localFtp.ReadFile(request.Parameters.PathLocalFtp, file)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

			return
		}

		//здесь нужно сохранять во временную папку
		// а затем отправлять в другой ftp серврер
		//после через полученныйй канал ответа, передавать ответ в NATS
	}
}

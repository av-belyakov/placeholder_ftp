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
		opts.Logger.Send("error", fmt.Sprintf("local FTP %s %s:%d", err.Error(), f, l-2))

		return
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfMainFtp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		opts.Logger.Send("error", fmt.Sprintf("main FTP %s %s:%d", err.Error(), f, l-2))

		return
	}

	request := RequestCopyFileFromFtpServer{}
	if err := json.Unmarshal(req.GetData(), &request); err != nil {
		_, f, l, _ := runtime.Caller(0)
		opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

		return
	}

	for _, fileName := range request.Parameters.Files {
		_, f, l, _ := runtime.Caller(0)
		countByteRead, err := localFtp.ReadFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: request.Parameters.PathLocalFtp,
				SrcFileName: fileName,
				DstFilePath: opts.TmpDir,
				DstFileName: fileName,
			})
		if err != nil {
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l+1))

			return
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully created", countByteRead, fileName))
		//********************************

		countByteWrite, err := mainFtp.WriteFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: opts.TmpDir,
				SrcFileName: fileName,
				DstFilePath: request.Parameters.PathMainFtp,
				DstFileName: fileName,
			})
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

			return
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully copied to FTP", countByteWrite, fileName))
		//********************************
	}
}

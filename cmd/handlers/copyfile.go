package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

// /HandlerCopyFile обработчик копирования файлов
func (opts FtpHandlerOptions) HandlerCopyFile(
	ctx context.Context,
	req ci.ChannelRequester) {

	result := NewResultRequestCopyFileFromFtpServer()
	result.SetRequestId(req.GetRequestId())

	// исходный ftp сервер
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfLocalFtp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		msgErr := fmt.Errorf("local FTP %s", err)
		opts.Logger.Send("error", fmt.Sprintf("%v %s:%d", msgErr, f, l-2))

		result.SetError(msgErr)

		return
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfMainFtp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		msgErr := fmt.Errorf("main FTP %s", err)
		opts.Logger.Send("error", fmt.Sprintf("%v %s:%d", msgErr, f, l-2))

		result.SetError(msgErr)

		return
	}

	request := RequestCopyFileFromFtpServer{}
	if err := json.Unmarshal(req.GetData(), &request); err != nil {
		_, f, l, _ := runtime.Caller(0)
		opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

		result.SetError(err)

		return
	}

	listProcessedFile := []commoninterfaces.FileInformationTransfer(nil)
	for _, fileName := range request.Parameters.Files {
		pf := NewProcessedFiles()
		pf.SetFileName(fileName)

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
			pf.Error = err
			listProcessedFile = append(listProcessedFile, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully created", countByteRead, fileName))
		//********************************

		_, f, l, _ = runtime.Caller(0)
		err = mainFtp.WriteFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: opts.TmpDir,
				SrcFileName: fileName,
				DstFilePath: request.Parameters.PathMainFtp,
				DstFileName: fileName,
			})
		if err != nil {
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l+1))
			pf.Error = err
			listProcessedFile = append(listProcessedFile, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("file '%s' has been successfully copied to FTP", fileName))
		//********************************

		_, f, l, _ = runtime.Caller(0)
		if err = os.Remove(path.Join(opts.TmpDir, fileName)); err != nil {
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l+1))
			pf.Error = err
		}

		pf.SizeBeforProcessing = countByteRead
		pf.SizeAfterProcessing = countByteRead

		listProcessedFile = append(listProcessedFile, pf)
	}

	result.SetData(listProcessedFile)

	ch := req.GetChanOutput()
	ch <- result
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

// /HandlerCopyFile обработчик копирования файлов
func (opts FtpHandlerOptions) HandlerCopyFile(ctx context.Context, req ci.ChannelRequester) {
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

	listProcessedLink := []commoninterfaces.LinkInformationTransfer(nil)
	for _, link := range request.Parameters.Links {
		pf := NewProcessedLink()
		pf.SetLinkOld(link)

		if ok := strings.HasPrefix(link, "ftp://"); !ok {
			err := fmt.Errorf("incorrect prefix in the link '%s'", link)

			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		result, err := supportingfunctions.LinkParse(link)
		if err != nil {
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		//чтение файла с ftp сервера источника
		countByteRead, err := localFtp.ReadFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: result.Path,
				SrcFileName: result.FileName,
				DstFilePath: opts.TmpDir,
				DstFileName: result.FileName,
			})
		if err != nil {
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())
			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully created", countByteRead, result.FileName))
		//********************************

		//формируем и устанавливаем ссылку по которой на MainFTP будет хранится файл
		u := &url.URL{
			Scheme: result.Scheme,
			Host:   opts.ConfMainFtp.GetHost(),
			Path:   path.Join(opts.PathResultDirMainFTP, result.FileName),
		}
		pf.SetLinkNew(u.String())

		//запись загрузка файла на ftp сервер назначения
		err = mainFtp.WriteFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: opts.TmpDir,
				SrcFileName: result.FileName,
				DstFilePath: opts.PathResultDirMainFTP,
				DstFileName: result.FileName,
			})
		if err != nil {
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())
			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("file '%s' has been successfully copied to FTP", result.FileName))
		//********************************

		var countByteDecode int
		if fi, err := os.Stat(path.Join(opts.TmpDir, result.FileName)); err == nil {
			countByteDecode = int(fi.Size())
		}

		//удаление временных файлов
		if err = deleteTmpFiles(opts.TmpDir, result.FileName, result.FileName); err != nil {
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())
			pf.SetError(err)
		}

		pf.SetSizeBeforProcessing(countByteRead)
		pf.SetSizeAfterProcessing(countByteDecode)

		listProcessedLink = append(listProcessedLink, pf)
	}

	result.SetData(listProcessedLink)

	req.GetChanOutput() <- result
}

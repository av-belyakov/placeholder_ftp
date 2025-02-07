package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
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
		msgErr := fmt.Errorf("local FTP %w (task id: '%s')", err, req.GetRequestId())
		opts.Logger.Send("error", supportingfunctions.CustomError(msgErr).Error())

		result.SetError(msgErr)

		return
	}

	opts.Logger.Send("info", fmt.Sprintf("successful connection to the LOCAL ftp server (task id: '%s')", req.GetRequestId()))

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfMainFtp)
	if err != nil {
		msgErr := fmt.Errorf("main FTP %w (task id: '%s')", err, req.GetRequestId())
		opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

		result.SetError(msgErr)

		return
	}

	request := RequestCopyFileFromFtpServer{}
	if err := json.Unmarshal(req.GetData(), &request); err != nil {
		err = fmt.Errorf("%w (task id: '%s')", err, req.GetRequestId())
		opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

		result.SetError(err)

		return
	}

	opts.Logger.Send("info", fmt.Sprintf("successful connection to the MAIN ftp server (task id: '%s')", req.GetRequestId()))

	listProcessedLink := []commoninterfaces.LinkInformationTransfer(nil)
	for _, link := range request.Parameters.Links {
		pf := NewProcessedLink()
		pf.SetLinkOld(link)

		if ok := strings.HasPrefix(link, "ftp://"); !ok {
			err := fmt.Errorf("incorrect prefix in the link '%s' (task id: '%s')", link, req.GetRequestId())
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		result, err := supportingfunctions.LinkParse(link)
		if err != nil {
			err = fmt.Errorf("%w (task id: '%s')", err, req.GetRequestId())
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
			err = fmt.Errorf("%w (task id: '%s')", err, req.GetRequestId())
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully created (task id: '%s')", countByteRead, result.FileName, req.GetRequestId()))

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
			err = fmt.Errorf("%w (task id: '%s')", err, req.GetRequestId())
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		opts.Logger.Send("info", fmt.Sprintf("the file '%s' was successfully decoded (task id: '%s')", result.FileName, req.GetRequestId()))

		var countByteDecode int
		if fi, err := os.Stat(path.Join(opts.TmpDir, result.FileName)); err == nil {
			countByteDecode = int(fi.Size())
		}

		//удаление временных файлов
		if err = deleteTmpFiles(opts.TmpDir, result.FileName, result.FileName); err != nil {
			err = fmt.Errorf("%w (task id: '%s')", err, req.GetRequestId())
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

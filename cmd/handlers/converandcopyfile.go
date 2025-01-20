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
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

// HandlerConvertAndCopyFile обработчик копирования и преобразования файлов
func (opts FtpHandlerOptions) HandlerConvertAndCopyFile(ctx context.Context, req commoninterfaces.ChannelRequester) {
	result := NewResultRequestCopyFileFromFtpServer()
	result.SetRequestId(req.GetRequestId())

	// исходный ftp сервер
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfLocalFtp)
	if err != nil {
		msgErr := fmt.Errorf("local FTP %w", err)
		opts.Logger.Send("error", supportingfunctions.CustomError(msgErr).Error())

		result.SetError(msgErr)
		req.GetChanOutput() <- result

		return
	}

	// ftp сервер назначения
	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(opts.ConfMainFtp)
	if err != nil {
		msgErr := fmt.Errorf("main FTP %w", err)
		opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

		result.SetError(msgErr)
		req.GetChanOutput() <- result

		return
	}

	request := RequestCopyFileFromFtpServer{}
	if err := json.Unmarshal(req.GetData(), &request); err != nil {
		opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())

		result.SetError(err)
		req.GetChanOutput() <- result

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

		suffTdp := strings.HasSuffix(link, ".tdp")
		suffPcap := strings.HasSuffix(link, ".pcap")
		if !suffTdp && !suffPcap {
			err := fmt.Errorf("incorrect suffix in the link '%s'", link)

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

		newFileName := result.FileName + ".txt"

		//формируем и устанавливаем ссылку по которой на MainFTP будет хранится файл
		u := &url.URL{
			Scheme: result.Scheme,
			Host:   opts.ConfMainFtp.GetHost(),
			Path:   path.Join(opts.PathResultDirMainFTP, newFileName),
		}
		pf.SetLinkNew(u.String())

		//декодирование и конвертация файла формата .pcap в текстовый вид
		if err = convertingNetworkTraffic(opts.TmpDir, result.FileName, newFileName, opts.Logger); err != nil {
			opts.Logger.Send("error", supportingfunctions.CustomError(err).Error())
			pf.SetError(err)
			listProcessedLink = append(listProcessedLink, pf)

			continue
		}

		//запись загрузка файла на ftp сервер назначения
		err = mainFtp.WriteFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: opts.TmpDir,
				SrcFileName: newFileName,
				DstFilePath: opts.PathResultDirMainFTP,
				DstFileName: newFileName,
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
		if fi, err := os.Stat(path.Join(opts.TmpDir, newFileName)); err == nil {
			countByteDecode = int(fi.Size())
		}

		//удаление временных файлов
		if err = deleteTmpFiles(opts.TmpDir, result.FileName, newFileName); err != nil {
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

func deleteTmpFiles(pathDir string, files ...string) error {
	for _, file := range files {
		if err := os.Remove(path.Join(pathDir, file)); err != nil {
			return err
		}
	}

	return nil
}

func convertingNetworkTraffic(filePath, rfn, wfn string, logging commoninterfaces.Logger) error {
	// для файла по которому выполняется декодирование пакетов
	readFile, err := os.Open(path.Join(filePath, rfn))
	if err != nil {
		return err
	}

	defer readFile.Close()

	// для файла в который выполняется запись информации полученной в результате декодирования
	writeFile, err := os.OpenFile(path.Join(filePath, wfn), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer writeFile.Close()

	err = supportingfunctions.NetworkTrafficDecoder(rfn, readFile, writeFile, logging)
	if err != nil {
		return err
	}

	return nil
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

// HandlerConvertAndCopyFile обработчик копирования и преобразования ффайлов
func (opts FtpHandlerOptions) HandlerConvertAndCopyFile(ctx context.Context, req commoninterfaces.ChannelRequester) {
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
		pf.SetFileNameOld(fileName)
		pf.SetFileNameNew(fmt.Sprintf("%s.txt", fileName))

		//чтение файла с ftp сервера источника
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
			pf.SetError(err)
			listProcessedFile = append(listProcessedFile, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("%d byte file '%s' has been successfully created", countByteRead, fileName))
		//********************************

		//декодирование и конвертация файла формата .pcap в текстовый вид

		if err = convertingNetworkTraffic(opts.TmpDir, pf.GetFileNameOld(), pf.GetFileNameNew(), opts.Logger); err != nil {
			_, f, l, _ = runtime.Caller(0)
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-1))
			pf.SetError(err)
			listProcessedFile = append(listProcessedFile, pf)

			continue
		}

		//запись загрузка файла на ftp сервер назначения
		err = mainFtp.WriteFile(
			ctx,
			wrappers.WrapperReadWriteFileOptions{
				SrcFilePath: opts.TmpDir,
				SrcFileName: pf.GetFileNameNew(),
				DstFilePath: request.Parameters.PathMainFtp,
				DstFileName: pf.GetFileNameNew(),
			})
		if err != nil {
			_, f, l, _ = runtime.Caller(0)
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-9))
			pf.SetError(err)
			listProcessedFile = append(listProcessedFile, pf)

			continue
		}

		//
		// это пока только для тестов
		//********************************
		opts.Logger.Send("info", fmt.Sprintf("file '%s' has been successfully copied to FTP", fileName))
		//********************************

		var countByteDecode int
		if fi, err := os.Stat(path.Join(opts.TmpDir, pf.GetFileNameNew())); err == nil {
			countByteDecode = int(fi.Size())
		}

		//удаление временных файлов
		if err = deleteTmpFiles(opts.TmpDir, pf.GetFileNameOld(), pf.GetFileNameNew()); err != nil {
			_, f, l, _ = runtime.Caller(0)
			opts.Logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-1))
			pf.SetError(err)
		}

		pf.SetSizeBeforProcessing(countByteRead)
		pf.SetSizeAfterProcessing(countByteDecode)

		listProcessedFile = append(listProcessedFile, pf)
	}

	result.SetData(listProcessedFile)

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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/av-belyakov/simplelogger"

	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/cmd/elasticsearchapi"
	"github.com/av-belyakov/placeholder_ftp/cmd/handlers"
	"github.com/av-belyakov/placeholder_ftp/cmd/messagebrokers/natsapi"
	"github.com/av-belyakov/placeholder_ftp/internal/appname"
	"github.com/av-belyakov/placeholder_ftp/internal/appversion"
	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
	"github.com/av-belyakov/placeholder_ftp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.New(rootPath, Conf_Dir)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	//******************************************************
	//********** инициализация модуля логирования **********
	loggingConf := confApp.GetSimpleLoggerPackage()
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, Root_Dir, getLoggerSettings(loggingConf))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %s", err.Error())
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	confDB := confApp.GetConfigWriteLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: confApp.NameRegionalObject,
	}); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-7))

		log.Println(err.Error())
	} else {
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//*****************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	zabbixConf := confApp.GetZabbixAPI()
	channelZabbix := make(chan ci.Messager)
	wzis := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: zabbixConf.NetworkPort,
		NetworkHost: zabbixConf.NetworkHost,
		ZabbixHost:  zabbixConf.ZabbixHost}

	eventTypes := []wrappers.EventType(nil)
	for _, v := range zabbixConf.EventTypes {
		eventTypes = append(eventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	wzis.EventTypes = eventTypes
	wrappers.WrappersZabbixInteraction(ctx, wzis, simpleLogger, channelZabbix)

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, simpleLogger, channelZabbix, logging.GetChan())

	//***************************************************
	//********** инициализация NATS API модуля **********
	confNatsSAPI := confApp.GetConfigNATS()
	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(confNatsSAPI.Host),
		natsapi.WithPort(confNatsSAPI.Port),
		natsapi.WithCacheTTL(confNatsSAPI.CacheTTL),
		natsapi.WithSubListenerCommand(confNatsSAPI.Subscriptions.ListenerCommand)}
	apiNats, err := natsapi.New(logging, confApp.GetNameRegionalObject(), natsOptsAPI...)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chNatsReqApi, err := apiNats.Start(ctx)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	confLocalFtp := confApp.GetConfigLocalFTP()
	confMainFtp := confApp.GetConfigMainFTP()

	//****************** проверка наличия доступа к FTP серверам ********************
	msgErr := "access initialization error"
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s LOCALFTP '%w'", msgErr, err)).Error())
		log.Fatalf("%s LOCALFTP '%s'\n", msgErr, err.Error())
	}

	//проверяем доступ к локальному ftp серверу
	if err := localFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s LOCALFTP '%w'", msgErr, err)).Error())
		log.Fatalf("%s LOCALFTP '%s'\n", msgErr, err.Error())
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s MAINFTP '%w'", msgErr, err)).Error())
		log.Fatalf("%s MAINFTP '%s'\n", msgErr, err.Error())
	}

	//проверяем доступ к удаленному ftp серверу
	if err = mainFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s MAINFTP '%w'", msgErr, err)).Error())
		log.Fatalf("%s MAINFTP '%s'\n", msgErr, err.Error())
	}
	//*******************************************************************************

	ftpho := handlers.FtpHandlerOptions{
		TmpDir:               Tmp_Files,
		PathResultDirMainFTP: confApp.MainFTPPathResultDirectory,
		ConfLocalFtp:         &confLocalFtp,
		ConfMainFtp:          &confMainFtp,
		Logger:               logging,
	}
	handlerList := map[string]func(ci.ChannelRequester){
		"copy_file": func(req ci.ChannelRequester) {
			ftpho.HandlerCopyFile(ctx, req)
		},
		"convert_and_copy_file": func(req ci.ChannelRequester) {
			ftpho.HandlerConvertAndCopyFile(ctx, req)
		},
	}

	//создание временной директории если ее нет
	if err := supportingfunctions.CreateDirectory(Root_Dir, Tmp_Files); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("error create tmp directory '%w'", err)).Error())
		log.Fatalf("error create tmp directory '%s'\n", err.Error())
	}

	appStatus := fmt.Sprintf("%vproduction%v", Ansi_Bright_Blue, Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_PHFTP_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", Ansi_Bright_Red, envValue, Ansi_Reset)
	}

	/*
	   	Делаю поочередно следующее:

	   0. Проверить что в ответном JSON код ответа был 200 при успешном выполнении задачи или
	   550 при не успешном.
	   Или может быть вообще убрать его?????? Потому что при если один файл был успешно обработан,
	   а другой нет то не понятно какой код ответа ставить.

	   1. В config добавляю параметр содержащий путь до директории, на MIIN_FTP, в котором хранятся файлы переведенные в формат txt.

	   2. Изменяю структуру запроса, со всей вытекающей логикой его обработки, на запрос вида:
	   ```
	   {
	     "task_id": "", //идентификатор задачи
	     "source": "", //наименование регионального объекта к которому был адресован запрос
	     "service": "", //имя сервиса-инициатора команды
	     "command": "convert_and_copy_file", //наименование команды
	     "parameters": {
	         "links": [
	             "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file.pcap",
	             "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file_http.pcap",
	             "..."
	             ] //список файлов которые необходимо обработать
	         }
	   }
	   ```

	   3. Изменяю структуру ответа, соответственно логику его формирующую тоже, на ответ вида:
	   ```
	   {
	     "request_id":"", //идентификатор задачи
	     "source": "", //наименование регионального объекта к которому был адресован запрос
	     "error": "", //содержит глобальные ошибки, такие как например, ошибка подключения к ftp серверу
	     "status_code": "", //код статуса выполнения задачи
	     "processed": [
	         {
	             "error": "" //ошибка возникшая при обработки файла
	             "size_befor_processing": int //размер файла до обработки
	             "size_after_processing": int //размер файла после обработки
	             "link_old": "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file.pcap",
	             "link_new": "ftp://ftp.cloud.gcm/traff/test_pcap_file.pcap.txt"
	         }
	     ]
	   }
	   ```

	   4. При обработки pcap результирующий txt файл не должен превышать размер в 50 Мб.
	*/

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetAppName(), appversion.GetAppVersion())
	fmt.Printf("\n%v%v%s.%v\n", Bold_Font, Ansi_Bright_Green, msg, Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", Underlining, Ansi_Bright_Green, appStatus, Ansi_Reset)
	fmt.Printf("%vLocal FTP server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%v  ip: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, confLocalFtp.Host, Ansi_Reset)
	fmt.Printf("%v  net port: %v%d%v\n", Ansi_Bright_Green, Ansi_Bright_Magenta, confLocalFtp.Port, Ansi_Reset)
	fmt.Printf("%vMain FTP server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%v  ip: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, confMainFtp.Host, Ansi_Reset)
	fmt.Printf("%v  net port: %v%d%v\n\n", Ansi_Bright_Green, Ansi_Bright_Magenta, confMainFtp.Port, Ansi_Reset)
	_ = simpleLogger.Write("info", strings.ToLower(msg))

	if err = router(ctx, handlerList, chNatsReqApi); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalln(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/av-belyakov/simplelogger"

	ci "github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/cmd/elasticsearchapi"
	"github.com/av-belyakov/placeholder_ftp/cmd/handlers"
	"github.com/av-belyakov/placeholder_ftp/cmd/messagebrokers/natsapi"
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

	// ********************************************************************************
	// ********************* инициализация модуля логирования *************************
	var listLog []simplelogger.OptionsManager
	for _, v := range confApp.GetSimpleLoggerPackage() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
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
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Println(err.Error())
	} else {
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//******************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	zabbixConf := confApp.GetZabbixAPI()
	chZabbix := make(chan ci.Messager)
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
	wrappers.WrappersZabbixInteraction(ctx, wzis, simpleLogger, chZabbix)

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	//******************************************************************
	//****************** инициализация NATS API модуля *****************
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
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s LOCALFTP '%w' (host:'%s' user: '%s')", msgErr, err, confLocalFtp.GetHost(), confLocalFtp.GetUsername())).Error())
		log.Fatalf("%s LOCALFTP '%s'\n", msgErr, err.Error())
	}

	//проверяем доступ к локальному ftp серверу
	if err := localFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s LOCALFTP '%w' (host:'%s' user: '%s')", msgErr, err, confLocalFtp.GetHost(), confLocalFtp.GetUsername())).Error())
		log.Fatalf("%s LOCALFTP '%s'\n", msgErr, err.Error())
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s MAINFTP '%w' (host:'%s' user: '%s')", msgErr, err, confMainFtp.GetHost(), confMainFtp.GetUsername())).Error())
		log.Fatalf("%s MAINFTP '%s'\n", msgErr, err.Error())
	}

	//проверяем доступ к удаленному ftp серверу
	if err = mainFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(fmt.Errorf("%s MAINFTP '%w' (host:'%s' user: '%s')", msgErr, err, confMainFtp.GetHost(), confMainFtp.GetUsername())).Error())
		log.Fatalf("%s MAINFTP '%s'\n", msgErr, err.Error())
	}
	//*******************************************************************************

	ftpho := handlers.FtpHandlerOptions{
		TmpDir:               Tmp_Files,
		PathResultDirMainFTP: confApp.MainFTPPathResultDirectory,
		MaxWritingFileLimit:  confApp.MaxWritingFileLimit,
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

	infoMsg := getInformationMessage(confApp.NameRegionalObject, confLocalFtp, confMainFtp)
	// вывод информационного сообщения при старте приложения
	_ = simpleLogger.Write("info", strings.ToLower(infoMsg))

	if err = router(ctx, handlerList, chNatsReqApi); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalln(err)
	}
}

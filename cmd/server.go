package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_ftp/cmd/handlers"
	"github.com/av-belyakov/placeholder_ftp/cmd/messagebrokers/natsapi"
	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
	"github.com/av-belyakov/placeholder_ftp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(ROOT_DIR)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.New(rootPath, CONF_DIR)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	//******************************************************
	//********** инициализация модуля логирования **********
	loggingConf := confApp.GetSimpleLoggerPackage()
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, ROOT_DIR, getLoggerSettings(loggingConf))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %s", err.Error())
	}

	//*****************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	zabbixConf := confApp.GetZabbixAPI()
	channelZabbix := make(chan commoninterfaces.Messager)
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

	if err := wrappers.WrappersZabbixInteraction(ctx, simpleLogger, wzis, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1), "error")
	}

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
	apiNats, err := natsapi.New(logging, natsOptsAPI...)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chNatsAPIReq, err := apiNats.Start(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	confLocalFtp := confApp.GetConfigLocalFTP()
	confMainFtp := confApp.GetConfigMainFTP()

	//****************** проверка наличия доступа к FTP серверам ********************
	msgErr := "access initialization error"
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("%s LOCALFTP '%s' %s:%d", msgErr, err.Error(), f, l-1), "error")
		log.Fatalf("%s LOCALFTP '%s': %s\n", msgErr, err.Error())
	}
	if err := localFtp.CheckConn(); err != nil {
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("%s LOCALFTP '%s' %s:%d", msgErr, err.Error(), f, l-1), "error")
		log.Fatalf("%s LOCALFTP '%s': %s\n", msgErr, err.Error())
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	_, f, l, _ = runtime.Caller(0)
	if err != nil {
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("%s MAINFTP '%s' %s:%d", msgErr, err.Error(), f, l-1), "error")
		log.Fatalf("%s MAINFTP '%s': %s\n", msgErr, err.Error())
	}
	if err = mainFtp.CheckConn(); err != nil {
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("%s MAINFTP '%s' %s:%d", msgErr, err.Error(), f, l-1), "error")
		log.Fatalf("%s MAINFTP '%s': %s\n", msgErr, err.Error())
	}
	//*******************************************************************************

	ftpho := handlers.FtpHandlerOptions{
		TmpDir:       TMP_FILES,
		ConfLocalFtp: &confLocalFtp,
		ConfMainFtp:  &confMainFtp,
		Logger:       logging,
	}
	handlers := map[string]func(commoninterfaces.ChannelRequester){
		"copy_file": func(req commoninterfaces.ChannelRequester) {
			ftpho.HandlerCopyFile(ctx, req)
		},
		"convert_and_copy_file": func(req commoninterfaces.ChannelRequester) {
			ftpho.HandlerConvertAndCopyFile(ctx, req)
		},
	}

	//создание временной директории если ее нет
	if err := supportingfunctions.CreateDirectory(ROOT_DIR, TMP_FILES); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf("error create tmp directory '%s' %s:%d", err.Error(), f, l-1), "error")
		log.Fatalf("error create tmp directory '%s'\n", err.Error())
	}

	router(ctx, handlers, chNatsAPIReq)
}

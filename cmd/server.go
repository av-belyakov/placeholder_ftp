package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
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
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1), "error")
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
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chNatsAPIReq, err := apiNats.Start(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	router(ctx, chNatsAPIReq)
}

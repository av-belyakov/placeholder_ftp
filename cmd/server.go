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
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chNatsReqApi, err := apiNats.Start(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	confLocalFtp := confApp.GetConfigLocalFTP()
	confMainFtp := confApp.GetConfigMainFTP()

	//****************** проверка наличия доступа к FTP серверам ********************
	msgErr := "access initialization error"
	localFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		_ = simpleLogger.Write("error", fmt.Sprintf("%s LOCALFTP '%s' %s:%d", msgErr, err.Error(), f, l-1))
		log.Fatalf("%s LOCALFTP '%s' %s:%d\n", msgErr, err.Error(), f, l-1)
	}

	//проверяем доступ к локальному ftp серверу
	if err := localFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", fmt.Sprintf("%s LOCALFTP '%s' %s:%d", msgErr, err.Error(), f, l-1))
		log.Fatalf("%s LOCALFTP '%s': %s:%d\n", msgErr, err.Error(), f, l-1)
	}

	mainFtp, err := wrappers.NewWrapperSimpleNetworkClient(&confLocalFtp)
	_, f, l, _ = runtime.Caller(0)
	if err != nil {
		_ = simpleLogger.Write("error", fmt.Sprintf("%s MAINFTP '%s' %s:%d", msgErr, err.Error(), f, l-1))
		log.Fatalf("%s MAINFTP '%s': %s:%d\n", msgErr, err.Error(), f, l-1)
	}

	//проверяем доступ к удаленному ftp серверу
	if err = mainFtp.CheckConn(); err != nil {
		_ = simpleLogger.Write("error", fmt.Sprintf("%s MAINFTP '%s' %s:%d", msgErr, err.Error(), f, l-1))
		log.Fatalf("%s MAINFTP '%s': %s:%d\n", msgErr, err.Error(), f, l-1)
	}
	//*******************************************************************************

	ftpho := handlers.FtpHandlerOptions{
		TmpDir:       Tmp_Files,
		ConfLocalFtp: &confLocalFtp,
		ConfMainFtp:  &confMainFtp,
		Logger:       logging,
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
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("error create tmp directory '%s' %s:%d", err.Error(), f, l-1))
		log.Fatalf("error create tmp directory '%s'\n", err.Error())
	}

	appStatus := fmt.Sprintf("%vproduction%v", Ansi_Bright_Blue, Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_PHFTP_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", Ansi_Bright_Red, envValue, Ansi_Reset)
	}

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
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1))
		log.Fatalln(err)
	}
}

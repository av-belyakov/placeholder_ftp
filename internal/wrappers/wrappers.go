package wrappers

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
	zabbixapicommunicator "github.com/av-belyakov/zabbixapicommunicator/cmd"
)

// WrappersZabbixInteraction обертка для взаимодействия с модулем zabbixapi
func WrappersZabbixInteraction(
	ctx context.Context,
	writerLoggingData commoninterfaces.WriterLoggingData,
	settings WrappersZabbixInteractionSettings,
	channelZabbix <-chan commoninterfaces.Messager) error {

	connTimeout := time.Duration(7 * time.Second)
	zc, err := zabbixapicommunicator.New(zabbixapicommunicator.SettingsZabbixConnection{
		Port:              settings.NetworkPort,
		Host:              settings.NetworkHost,
		NetProto:          "tcp",
		ZabbixHost:        settings.ZabbixHost,
		ConnectionTimeout: &connTimeout,
	})
	if err != nil {
		return err
	}

	et := make([]zabbixapicommunicator.EventType, len(settings.EventTypes))
	for _, v := range settings.EventTypes {
		et = append(et, zabbixapicommunicator.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  zabbixapicommunicator.Handshake(v.Handshake),
		})
	}

	recipient := make(chan zabbixapicommunicator.Messager)
	if err = zc.Start(ctx, et, recipient); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-channelZabbix:
				newMessageSettings := &zabbixapicommunicator.MessageSettings{}
				newMessageSettings.SetType(msg.GetType())
				newMessageSettings.SetMessage(msg.GetMessage())

				recipient <- newMessageSettings
			}
		}
	}()

	go func() {
		for err := range zc.GetChanErr() {
			_, f, l, _ := runtime.Caller(0)
			writerLoggingData.WriteLoggingData(fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1), "error")
		}
	}()

	return nil
}

// NewWrapperSimpleNetworkClient формирует обертку для взаимодействия с FTP клиентами
func NewWrapperSimpleNetworkClient(settings commoninterfaces.SimpleNetworkConsumer) (*WrapperSimplyNetworkClient, error) {
	netClient := &WrapperSimplyNetworkClient{}

	if settings.GetHost() == "" {
		return netClient, fmt.Errorf("the value 'Host' should not be empty")
	}
	netClient.setHost(settings.GetHost())

	if settings.GetPort() == 0 {
		return netClient, fmt.Errorf("the value 'Port' should not be equal '0'")
	}
	netClient.setPort(settings.GetPort())

	if settings.GetUsername() == "" {
		return netClient, fmt.Errorf("the value 'Username' should not be empty")
	}
	netClient.setUsername(settings.GetUsername())

	if settings.GetPasswd() == "" {
		return netClient, fmt.Errorf("the value 'Passwd' should not be empty")
	}
	netClient.setPasswd(settings.GetPasswd())

	return netClient, nil
}

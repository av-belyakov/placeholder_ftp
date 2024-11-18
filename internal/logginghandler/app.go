package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/zabbixapicommunicator"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// LoggingHandler обработчик и распределитель логов
func LoggingHandler(
	ctx context.Context,
	writerLoggingData commoninterfaces.WriterLoggingData,
	channelZabbix chan<- commoninterfaces.Messager,
	logging <-chan commoninterfaces.Messager) {

	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-logging:
			//**********************************************************************
			//здесь так же может быть вывод в консоль, с индикацией цветов соответствующих
			//определенному типу сообщений но для этого надо включить вывод на stdout
			//в конфигурационном фале
			_ = writerLoggingData.WriteLoggingData(msg.GetMessage(), msg.GetType())

			if msg.GetType() == "error" || msg.GetType() == "warning" {
				channelZabbix <- &zabbixapicommunicator.MessageSettings{
					EventType: "error",
					Message:   fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()),
				}
			}

			if msg.GetType() == "info" {
				channelZabbix <- &zabbixapicommunicator.MessageSettings{
					EventType: "info",
					Message:   msg.GetMessage(),
				}
			}
		}
	}
}

package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// LoggingHandler обработчик и распределитель логов
func LoggingHandler(
	ctx context.Context,
	writerLoggingData commoninterfaces.WriterLoggingData,
	chanSystemMonitoring chan<- commoninterfaces.Messager,
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
			_ = writerLoggingData.Write(msg.GetType(), msg.GetMessage())

			if msg.GetType() == "error" || msg.GetType() == "warning" {
				msg := NewMessageLogging()
				msg.SetType("error")
				msg.SetMessage(fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()))

				chanSystemMonitoring <- msg
			}

			if msg.GetType() == "info" {
				msg := NewMessageLogging()
				msg.SetType("info")
				msg.SetMessage(msg.GetMessage())

				chanSystemMonitoring <- msg
			}
		}
	}
}

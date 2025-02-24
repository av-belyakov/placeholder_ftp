package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// New конструктор обработчиа логов (это просто мост соединяющий несколько сервисов логирования)
func New(writer commoninterfaces.WriterLoggingData, chSysMonit chan<- commoninterfaces.Messager) *LoggingChan {
	return &LoggingChan{
		dataWriter:           writer,
		chanSystemMonitoring: chSysMonit,
		chanLogging:          make(chan commoninterfaces.Messager),
	}
}

// Start обработчик и распределитель логов
func (lc *LoggingChan) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-lc.chanLogging:
				//****************************************************************************
				//здесь так же может быть вывод в консоль, с индикацией цветов соответствующих
				//определенному типу сообщений но для этого надо включить вывод на stdout
				//в конфигурационном фале
				lc.dataWriter.Write(msg.GetType(), msg.GetMessage())

				if msg.GetType() == "error" || msg.GetType() == "warning" {
					msg := NewMessageLogging()
					msg.SetType("error")
					msg.SetMessage(fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()))

					lc.chanSystemMonitoring <- msg
				}

				if msg.GetType() == "info" {
					msg := NewMessageLogging()
					msg.SetType("info")
					msg.SetMessage(msg.GetMessage())

					lc.chanSystemMonitoring <- msg
				}
			}
		}
	}()
}

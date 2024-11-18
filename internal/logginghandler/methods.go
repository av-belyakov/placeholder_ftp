package logginghandler

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

func New() *LoggingChan {
	return &LoggingChan{
		logChan: make(chan commoninterfaces.Messager),
	}
}

func (l *LoggingChan) GetChan() <-chan commoninterfaces.Messager {
	return l.logChan
}

func (l *LoggingChan) Send(msgType, message string) {
	ms := NewMessageLogging()
	ms.SetType(msgType)
	ms.SetMessage(message)

	l.logChan <- ms
}

func (l *LoggingChan) Close() {
	close(l.logChan)
}

func NewMessageLogging() *MessageLogging {
	return &MessageLogging{}
}

// GetMessage возвращает сообщение
func (ml *MessageLogging) GetMessage() string {
	return ml.Message
}

// SetMessage устанавливает сообщение
func (ml *MessageLogging) SetMessage(v string) {
	ml.Message = v
}

// GetType возвращает тип сообщения
func (ml *MessageLogging) GetType() string {
	return ml.Type
}

// SetType устанавливает тип сообщения
func (ml *MessageLogging) SetType(v string) {
	ml.Type = v
}

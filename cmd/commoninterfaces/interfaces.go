package commoninterfaces

//************ каналы *************

type ChannelResponser interface {
	RequestIdHandler
	GetStatusCode() int
	SetStatusCode(int)
	GetError() error
	SetError(error)
	GetData() []byte
	SetData([]byte)
}

type ChannelRequester interface {
	RequestIdHandler
	GetData() interface{}
	SetData(interface{})
	GetChanOutput() chan ChannelResponser
	SetChanOutput(chan ChannelResponser)
}

type RequestIdHandler interface {
	GetRequestId() string
	SetRequestId(string)
}

//************** логирование ***************

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

type WriterLoggingData interface {
	WriteLoggingData(str, typeLogFile string) bool
}

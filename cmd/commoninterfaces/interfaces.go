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

//************** простой сетевой клиент ***************

type SimpleNetworkConsumer interface {
	HostHandler
	PortHandler
	UsernameHandler
	PasswdHandler
}

type HostHandler interface {
	GetHost() string
	SetHost(v string)
}

type PortHandler interface {
	GetPort() int
	SetPort(v int)
}

type UsernameHandler interface {
	GetUsername() string
	SetUsername(v string)
}

type PasswdHandler interface {
	GetPasswd() string
	SetPasswd(v string)
}

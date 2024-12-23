package commoninterfaces

//************ каналы *************

type ChannelResponser interface {
	RequestIdHandler
	GetError() error
	SetError(error)
	GetData() []FileInformationTransfer
	SetData([]FileInformationTransfer)
}

type ChannelRequester interface {
	RequestIdHandler
	GetCommand() string
	SetCommand(v string)
	GetOrder() string
	SetOrder(v string)
	GetData() []byte
	SetData([]byte)
	GetChanOutput() chan ChannelResponser
	SetChanOutput(chan ChannelResponser)
}

type FileInformationTransfer interface {
	ErrorHandler
	GetFileName() string
	SetFileName(v string)
	GetSizeBeforProcessing() int
	SetSizeBeforProcessing(int)
	GetSizeAfterProcessing() int
	SetSizeAfterProcessing(int)
}

type RequestIdHandler interface {
	GetRequestId() string
	SetRequestId(string)
}

type ErrorHandler interface {
	GetError() error
	SetError(error)
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

package wrappers

// WrappersZabbixInteractionSettings настройки для обертки взаимодействия с модулем zabbixapi
type WrappersZabbixInteractionSettings struct {
	NetworkPort int         //сетевой порт
	NetworkHost string      //ip адрес или доменное имя
	ZabbixHost  string      //zabbix host
	EventTypes  []EventType //типы событий
}

type EventType struct {
	IsTransmit bool
	EventType  string
	ZabbixKey  string
	Handshake  Handshake
}

type Handshake struct {
	TimeInterval int
	Message      string
}

// WrapperSimplyNetworkClient обертка с настройками для взаимодействия с простым сетевым клиентом
type WrapperSimplyNetworkClient struct {
	port     int
	host     string
	username string
	passwd   string
}

// WrapperReadWriteFileOptions опции для обертки методов чтения или записи файла
type WrapperReadWriteFileOptions struct {
	SrcFilePath string
	SrcFileName string
	DstFilePath string
	DstFileName string
}

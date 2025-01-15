package wrappers

// WrappersZabbixInteractionSettings настройки для обертки взаимодействия с модулем zabbixapi
type WrappersZabbixInteractionSettings struct {
	EventTypes  []EventType //типы событий
	NetworkHost string      //ip адрес или доменное имя
	ZabbixHost  string      //zabbix host
	NetworkPort int         //сетевой порт
}

type EventType struct {
	EventType  string
	ZabbixKey  string
	Handshake  Handshake
	IsTransmit bool
}

type Handshake struct {
	Message      string
	TimeInterval int
}

// WrapperSimplyNetworkClient обертка с настройками для взаимодействия с простым сетевым клиентом
type WrapperSimplyNetworkClient struct {
	host     string
	username string
	passwd   string
	port     int
}

// WrapperReadWriteFileOptions опции для обертки методов чтения или записи файла
type WrapperReadWriteFileOptions struct {
	SrcFilePath string
	SrcFileName string
	DstFilePath string
	DstFileName string
}

package confighandler

// AppConfig настройки приложения
type AppConfig struct {
	Information
	Logging    ConfigLogs
	NATS       ConfigNATS
	LocalFTP   ConfigFtp
	MainFTP    ConfigFtp
	WriteLogDB ConfigWriteLogDB
}

// Information информация о приложении
type Information struct {
	FileName           string `validate:"required" yaml:"filename"`
	NameRegionalObject string `validate:"required" yaml:"name_regional_object"`
}

// ConfigLogs настройки логирования
type ConfigLogs struct {
	ZabbixAPI           ZabbixOptions
	SimpleLoggerPackage []LoggerOption
}

type Logs struct {
	Logging []LoggerOption
}

type LoggerOption struct {
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
}

type ZabbixSet struct {
	Zabbix ZabbixOptions
}

type ZabbixOptions struct {
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	EventTypes  []EventType `yaml:"eventType"`
}

type EventType struct {
	IsTransmit bool      `yaml:"isTransmit"`
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
}

type Handshake struct {
	TimeInterval int    `yaml:"timeInterval"`
	Message      string `validate:"required" yaml:"message"`
}

// ConfigNATS настройки NATS
type ConfigNATS struct {
	Port int `validate:"gt=0,lte=65535" yaml:"port"`
	//сетевой порт
	CacheTTL int `validate:"gt=10,lte=86400" yaml:"cacheTtl"`
	//время жизни кеша
	Host string `validate:"required" yaml:"host"`
	//ip адрес или доменное имя
	Prefix string `yaml:"prefix"`
	//префикс
	Subscriptions SubscriptionsNATS `yaml:"subscriptions"`
	//список подписок
}

type SubscriptionsNATS struct {
	ListenerCommand string `validate:"required" yaml:"listener_command"`
}

type ConfigFtp struct {
	Port     int    `validate:"required" yaml:"port"`
	Host     string `validate:"required" yaml:"host"`
	Username string `validate:"required" yaml:"username"`
	Passwd   string `validate:"required" yaml:"passwd"`
}

type ConfigWriteLogDB struct {
	Port          int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	NameDB        string `yaml:"namedb"`
	StorageNameDB string `yaml:"storageNamedb"`
}

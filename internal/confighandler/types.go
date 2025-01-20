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
	FileName                   string `validate:"required" yaml:"filename"`
	NameRegionalObject         string `validate:"required" yaml:"name_regional_object"`
	MainFTPPathResultDirectory string `validate:"required" yaml:"main_ftp_path_result_directory"`
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
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
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
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
	IsTransmit bool      `yaml:"isTransmit"`
}

type Handshake struct {
	Message      string `validate:"required" yaml:"message"`
	TimeInterval int    `yaml:"timeInterval"`
}

// ConfigNATS настройки NATS
type ConfigNATS struct {
	Subscriptions SubscriptionsNATS `yaml:"subscriptions"`                        //список подписок
	Prefix        string            `yaml:"prefix"`                               //префикс
	Host          string            `validate:"required" yaml:"host"`             //ip адрес или доменное имя
	Port          int               `validate:"gt=0,lte=65535" yaml:"port"`       //сетевой порт
	CacheTTL      int               `validate:"gt=10,lte=86400" yaml:"cache_ttl"` //время жизни кеша
}

type SubscriptionsNATS struct {
	ListenerCommand string `validate:"required" yaml:"listener_command"`
}

type ConfigFtp struct {
	Host     string `validate:"required" yaml:"host"`
	Username string `validate:"required" yaml:"username"`
	Passwd   string `validate:"required" yaml:"passwd"`
	Port     int    `validate:"required" yaml:"port"`
}

type ConfigWriteLogDB struct {
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	NameDB        string `yaml:"name_db"`
	StorageNameDB string `yaml:"storage_name_db"`
	Port          int    `validate:"gt=0,lte=65535" yaml:"port"`
}

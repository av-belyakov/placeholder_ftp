package confighandler

// GetFileName имя конфигурационного файла
func (c AppConfig) GetFileName() string {
	return c.Information.FileName
}

// GetNameRegionalObject имя регионального объекта
func (c AppConfig) GetNameRegionalObject() string {
	return c.Information.NameRegionalObject
}

// GetMainFTPPathResultDirectory путь на MainFTP сервере для хранения файлов
func (c AppConfig) GetMainFTPPathResultDirectory() string {
	return c.Information.MainFTPPathResultDirectory
}

// GetSimpleLoggerPackage настройки пакета 'simplelogger'
func (c AppConfig) GetSimpleLoggerPackage() []LoggerOption {
	return c.Logging.SimpleLoggerPackage
}

// GetZabbixAPI настройки API Zabbix
func (c AppConfig) GetZabbixAPI() ZabbixOptions {
	return c.Logging.ZabbixAPI
}

// GetConfigNATS настройки для взаимодействия с брокером сообщений NATS
func (c AppConfig) GetConfigNATS() ConfigNATS {
	return c.NATS
}

// GetConfigLocalFTP настройки локального FTP серврера
func (c AppConfig) GetConfigLocalFTP() ConfigFtp {
	return c.LocalFTP
}

// GetConfigMainFTP настройки FTP сервера агрегатора файлов
func (c AppConfig) GetConfigMainFTP() ConfigFtp {
	return c.MainFTP
}

func (conf *ConfigFtp) GetHost() string {
	return conf.Host
}

func (conf *ConfigFtp) SetHost(v string) {
	conf.Host = v
}

func (conf *ConfigFtp) GetPort() int {
	return conf.Port
}

func (conf *ConfigFtp) SetPort(v int) {
	conf.Port = v
}

func (conf *ConfigFtp) GetUsername() string {
	return conf.Username
}

func (conf *ConfigFtp) SetUsername(v string) {
	conf.Username = v
}

func (conf *ConfigFtp) GetPasswd() string {
	return conf.Passwd
}

func (conf *ConfigFtp) SetPasswd(v string) {
	conf.Passwd = v
}

// GetConfigWriteLogDB настройки доступа к БД в которую осуществляется запись логов
func (c AppConfig) GetConfigWriteLogDB() ConfigWriteLogDB {
	return c.WriteLogDB
}

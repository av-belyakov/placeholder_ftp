package confighandler

// GetFileName имя конфигурационного файла
func (c AppConfig) GetFileName() string {
	return c.Information.FileName
}

// GetNameRegionalObject имя регионального объекта
func (c AppConfig) GetNameRegionalObject() string {
	return c.Information.NameRegionalObject
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

// GetConfigMainFTP настройки FTP серврера агрегатора файлов
func (c AppConfig) GetConfigMainFTP() ConfigFtp {
	return c.MainFTP
}

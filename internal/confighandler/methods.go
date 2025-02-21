package confighandler

import "errors"

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

// GetMaxWritingFileLimit ограничение максимального размера файла, в мегабайтах, который будет передаваться на ftp MainFTP
func (c AppConfig) GetMaxWritingFileLimit() int {
	return c.Information.MaxWritingFileLimit
}

// GetSimpleLoggerPackage настройки пакета 'simplelogger'
func (c AppConfig) GetSimpleLoggerPackage() []*LoggerOption {
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

// SetNameMessageType наименование тпа логирования
func (l *LoggerOption) SetNameMessageType(v string) error {
	if v == "" {
		return errors.New("the value 'MsgTypeName' must not be empty")
	}

	return nil
}

// SetMaxLogFileSize максимальный размер файла для логирования
func (l *LoggerOption) SetMaxLogFileSize(v int) error {
	if v < 1000 {
		return errors.New("the value 'MaxFileSize' must not be less than 1000")
	}

	return nil
}

// SetPathDirectory путь к директории логирования
func (l *LoggerOption) SetPathDirectory(v string) error {
	if v == "" {
		return errors.New("the value 'PathDirectory' must not be empty")
	}

	return nil
}

// SetWritingStdout запись логов на вывод stdout
func (l *LoggerOption) SetWritingStdout(v bool) {
	l.WritingStdout = v
}

// SetWritingFile запись логов в файл
func (l *LoggerOption) SetWritingFile(v bool) {
	l.WritingFile = v
}

// SetWritingDB запись логов  в БД
func (l *LoggerOption) SetWritingDB(v bool) {
	l.WritingDB = v
}

// GetNameMessageType наименование тпа логирования
func (l *LoggerOption) GetNameMessageType() string {
	return l.MsgTypeName
}

// GetMaxLogFileSize максимальный размер файла для логирования
func (l *LoggerOption) GetMaxLogFileSize() int {
	return l.MaxFileSize
}

// GetPathDirectory путь к директории логирования
func (l *LoggerOption) GetPathDirectory() string {
	return l.PathDirectory
}

// GetWritingStdout запись логов на вывод stdout
func (l *LoggerOption) GetWritingStdout() bool {
	return l.WritingStdout
}

// GetWritingFile запись логов в файл
func (l *LoggerOption) GetWritingFile() bool {
	return l.WritingFile
}

// GetWritingDB запись логов  в БД
func (l *LoggerOption) GetWritingDB() bool {
	return l.WritingDB
}

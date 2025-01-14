package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/placeholder_ftp/internal/supportingfunctions"
)

func New(rootDir, confDir string) (*AppConfig, error) {
	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_PHFTP_MAIN": "",

			//Подключение к NATS
			"GO_PHFTP_NPREFIX":             "",
			"GO_PHFTP_NHOST":               "",
			"GO_PHFTP_NPORT":               "",
			"GO_PHFTP_NCACHETTL":           "",
			"GO_PHFTP_NSUBLISTENERCOMMAND": "",

			//Подключение к локальному FTP серверу
			"GO_PHFTP_LOCALFTP_HOST":     "",
			"GO_PHFTP_LOCALFTP_PORT":     "",
			"GO_PHFTP_LOCALFTP_USERNAME": "",
			"GO_PHFTP_LOCALFTP_PASSWD":   "",

			//Подключение к FTP серверу агрегатору
			"GO_PHFTP_MAINFTP_HOST":     "",
			"GO_PHFTP_MAINFTP_PORT":     "",
			"GO_PHFTP_MAINFTP_USERNAME": "",
			"GO_PHFTP_MAINFTP_PASSWD":   "",

			//Подключение к БД в которую будут записыватся логи
			"GO_PHFTP_DBWLOGHOST":        "",
			"GO_PHFTP_DBWLOGPORT":        "",
			"GO_PHFTP_DBWLOGNAME":        "",
			"GO_PHFTP_DBWLOGSTORAGENAME": "",
			"GO_PHFTP_DBWLOGUSER":        "",
			"GO_PHFTP_DBWLOGPASSWD":      "",
		}
	)

	conf := AppConfig{}
	validate = validator.New(validator.WithRequiredStructEnabled())

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return &conf, err
	}

	confPath := path.Join(rootPath, confDir)

	list, err := os.ReadDir(confPath)
	if err != nil {
		return &conf, err
	}

	fileNameCommon, err := getFileName("config.yaml", confPath, list)
	if err != nil {
		return &conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon, &conf); err != nil {
		return &conf, err
	}

	var fn string
	if envList["GO_PHFTP_MAIN"] == "test" {
		fn, err = getFileName("config_test.yaml", confPath, list)
	} else if envList["GO_PHFTP_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yaml", confPath, list)
	} else {
		fn, err = getFileName("config_prod.yaml", confPath, list)
	}
	if err != nil {
		return &conf, err
	}

	if err := setSpecial(fn, &conf); err != nil {
		return &conf, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_PHFTP_NPREFIX"] != "" {
		conf.NATS.Prefix = envList["GO_PHFTP_NPREFIX"]
	}
	if envList["GO_PHFTP_NHOST"] != "" {
		conf.NATS.Host = envList["GO_PHFTP_NHOST"]
	}
	if envList["GO_PHFTP_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHFTP_NPORT"]); err == nil {
			conf.NATS.Port = p
		}
	}
	if envList["GO_PHFTP_NCACHETTL"] != "" {
		if v, err := strconv.Atoi(envList["GO_PHFTP_NCACHETTL"]); err == nil {
			conf.NATS.CacheTTL = v
		}
	}

	if envList["GO_PHFTP_NSUBLISTENERCOMMAND"] != "" {
		conf.NATS.Subscriptions.ListenerCommand = envList["GO_PHFTP_NSUBLISTENERCOMMAND"]
	}

	//Настройки локального FTP сервера
	if envList["GO_PHFTP_LOCALFTP_HOST"] != "" {
		conf.LocalFTP.Host = envList["GO_PHFTP_LOCALFTP_HOST"]
	}
	if envList["GO_PHFTP_LOCALFTP_PORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHFTP_LOCALFTP_PORT"]); err == nil {
			conf.LocalFTP.Port = p
		}
	}
	if envList["GO_PHFTP_LOCALFTP_USERNAME"] != "" {
		conf.LocalFTP.Username = envList["GO_PHFTP_LOCALFTP_USERNAME"]
	}
	if envList["GO_PHFTP_LOCALFTP_PASSWD"] != "" {
		conf.LocalFTP.Passwd = envList["GO_PHFTP_LOCALFTP_PASSWD"]
	}

	//Настройки FTP сервера агрегатора
	if envList["GO_PHFTP_MAINFTP_HOST"] != "" {
		conf.MainFTP.Host = envList["GO_PHFTP_MAINFTP_HOST"]
	}
	if envList["GO_PHFTP_MAINFTP_PORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHFTP_MAINFTP_PORT"]); err == nil {
			conf.MainFTP.Port = p
		}
	}
	if envList["GO_PHFTP_MAINFTP_USERNAME"] != "" {
		conf.MainFTP.Username = envList["GO_PHFTP_MAINFTP_USERNAME"]
	}
	if envList["GO_PHFTP_MAINFTP_PASSWD"] != "" {
		conf.MainFTP.Passwd = envList["GO_PHFTP_MAINFTP_PASSWD"]
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_PHFTP_DBWLOGHOST"] != "" {
		conf.WriteLogDB.Host = envList["GO_PHFTP_DBWLOGHOST"]
	}
	if envList["GO_PHFTP_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHFTP_DBWLOGPORT"]); err == nil {
			conf.WriteLogDB.Port = p
		}
	}
	if envList["GO_PHFTP_DBWLOGNAME"] != "" {
		conf.WriteLogDB.NameDB = envList["GO_PHFTP_DBWLOGNAME"]
	}
	if envList["GO_PHFTP_DBWLOGUSER"] != "" {
		conf.WriteLogDB.User = envList["GO_PHFTP_DBWLOGUSER"]
	}
	if envList["GO_PHFTP_DBWLOGPASSWD"] != "" {
		conf.WriteLogDB.Passwd = envList["GO_PHFTP_DBWLOGPASSWD"]
	}
	if envList["GO_PHFTP_DBWLOGSTORAGENAME"] != "" {
		conf.WriteLogDB.StorageNameDB = envList["GO_PHFTP_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return &conf, err
	}

	return &conf, nil
}

func getFileName(sf, confPath string, lfs []fs.DirEntry) (string, error) {
	for _, v := range lfs {
		if v.Name() == sf && !v.IsDir() {
			return path.Join(confPath, v.Name()), nil
		}
	}

	return "", fmt.Errorf("file '%s' is not found", sf)
}

func setCommonSettings(filename string, conf *AppConfig) error {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	setLoggingSettings := func() error {
		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			conf.Logging.SimpleLoggerPackage = ls.Logging
		}

		return nil
	}

	setZabbixAPISettings := func() error {
		z := ZabbixSet{}
		if ok := viper.IsSet("ZABBIX"); ok {
			if err := viper.GetViper().Unmarshal(&z); err != nil {
				return err
			}

			np := 10051
			if z.Zabbix.NetworkPort != 0 && z.Zabbix.NetworkPort < 65536 {
				np = z.Zabbix.NetworkPort
			}

			conf.Logging.ZabbixAPI = ZabbixOptions{
				NetworkPort: np,
				NetworkHost: z.Zabbix.NetworkHost,
				ZabbixHost:  z.Zabbix.ZabbixHost,
				EventTypes:  z.Zabbix.EventTypes,
			}
		}

		return nil
	}

	if err := setLoggingSettings(); err != nil {
		return err
	}

	if err := setZabbixAPISettings(); err != nil {
		return err
	}

	return nil
}

func setSpecial(filename string, conf *AppConfig) error {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if viper.IsSet("COMMONINFO.file_name") {
		conf.FileName = viper.GetString("COMMONINFO.file_name")
	}

	if viper.IsSet("COMMONINFO.name_regional_object") {
		conf.NameRegionalObject = viper.GetString("COMMONINFO.name_regional_object")
	}

	//Настройки для модуля подключения к NATS
	if viper.IsSet("NATS.prefix") {
		conf.NATS.Prefix = viper.GetString("NATS.prefix")
	}
	if viper.IsSet("NATS.host") {
		conf.NATS.Host = viper.GetString("NATS.host")
	}
	if viper.IsSet("NATS.port") {
		conf.NATS.Port = viper.GetInt("NATS.port")
	}
	if viper.IsSet("NATS.cacheTtl") {
		conf.NATS.CacheTTL = viper.GetInt("NATS.cacheTtl")
	}

	if viper.IsSet("NATS.subscriptions.listener_command") {
		conf.NATS.Subscriptions.ListenerCommand = viper.GetString("NATS.subscriptions.listener_command")
	}

	//Настройки локального FTP сервера
	if viper.IsSet("LOCALFTP.host") {
		conf.LocalFTP.Host = viper.GetString("LOCALFTP.host")
	}
	if viper.IsSet("LOCALFTP.port") {
		conf.LocalFTP.Port = viper.GetInt("LOCALFTP.port")
	}
	if viper.IsSet("LOCALFTP.username") {
		conf.LocalFTP.Username = viper.GetString("LOCALFTP.username")
	}

	//Настройки FTP сервера агрегатора
	if viper.IsSet("MAINFTP.host") {
		conf.MainFTP.Host = viper.GetString("MAINFTP.host")
	}
	if viper.IsSet("MAINFTP.port") {
		conf.MainFTP.Port = viper.GetInt("MAINFTP.port")
	}
	if viper.IsSet("MAINFTP.username") {
		conf.MainFTP.Username = viper.GetString("MAINFTP.username")
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if viper.IsSet("DATABASEWRITELOG.host") {
		conf.WriteLogDB.Host = viper.GetString("DATABASEWRITELOG.host")
	}
	if viper.IsSet("DATABASEWRITELOG.port") {
		conf.WriteLogDB.Port = viper.GetInt("DATABASEWRITELOG.port")
	}
	if viper.IsSet("DATABASEWRITELOG.user") {
		conf.WriteLogDB.User = viper.GetString("DATABASEWRITELOG.user")
	}
	if viper.IsSet("DATABASEWRITELOG.namedb") {
		conf.WriteLogDB.NameDB = viper.GetString("DATABASEWRITELOG.namedb")
	}
	if viper.IsSet("DATABASEWRITELOG.storageNamedb") {
		conf.WriteLogDB.StorageNameDB = viper.GetString("DATABASEWRITELOG.storageNamedb")
	}

	return nil
}

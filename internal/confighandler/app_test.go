package confighandler_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
)

const (
	CONF_DIR_NAME = "config"

	LFTP_PASSWD = "passwd_local_ftp"
	MFTP_PASSWD = "passwd_main_ftp"
)

func TestConfigHandler(t *testing.T) {
	os.Setenv("GO_PHFTP_LOCALFTP_PASSWD", LFTP_PASSWD)
	os.Setenv("GO_PHFTP_MAINFTP_PASSWD", MFTP_PASSWD)
	defer func() {
		os.Unsetenv("GO_PHFTP_LOCALFTP_PASSWD")
		os.Unsetenv("GO_PHFTP_MAINFTP_PASSWD")
	}()

	//для config_prod
	conf, err := confighandler.New("placeholder_ftp", CONF_DIR_NAME)
	assert.NoError(t, err)
	assert.Equal(t, conf.GetFileName(), "config_prod")
	assert.Greater(t, len(conf.GetSimpleLoggerPackage()), 0)

	//чтение тестового файла config_test
	os.Setenv("GO_PHFTP_MAIN", "test")

	conf, err = confighandler.New("placeholder_ftp", CONF_DIR_NAME)
	assert.NoError(t, err)

	assert.Equal(t, conf.GetFileName(), "config_test")
	assert.Equal(t, conf.GetNameRegionalObject(), "testrcm")

	//Подключение к NATS
	confNats := conf.GetConfigNATS()
	assert.Equal(t, confNats.Host, "nats.cloud.gcm")
	assert.Equal(t, confNats.Port, 4222)
	assert.Equal(t, confNats.Prefix, "test")
	assert.Equal(t, confNats.CacheTTL, 3600)
	assert.Equal(t, confNats.Subscriptions.ListenerCommand, "phftp.commands.test")

	confLocalFTP := conf.GetConfigLocalFTP()
	assert.Equal(t, confLocalFTP.Host, "127.0.0.1")
	assert.Equal(t, confLocalFTP.Username, "userlocalftp.test")
	assert.Equal(t, confLocalFTP.Passwd, LFTP_PASSWD)

	confMainFTP := conf.GetConfigMainFTP()
	assert.Equal(t, confMainFTP.Host, "127.0.0.1")
	assert.Equal(t, confMainFTP.Username, "usermainftp.test")
	assert.Equal(t, confMainFTP.Passwd, MFTP_PASSWD)

	os.Setenv("GO_PHFTP_NPREFIX", "new_prefix")
	os.Setenv("GO_PHFTP_NHOST", "localhost")
	os.Setenv("GO_PHFTP_NPORT", "3344")
	os.Setenv("GO_PHFTP_NCACHETTL", "4800")
	os.Setenv("GO_PHFTP_NSUBLISTENERCOMMAND", "phftp.commands.test")
	defer func() {
		os.Unsetenv("GO_PHFTP_NPREFIX")
		os.Unsetenv("GO_PHFTP_NHOST")
		os.Unsetenv("GO_PHFTP_NPORT")
		os.Unsetenv("GO_PHFTP_NCACHETTL")
		os.Unsetenv("GO_PHFTP_NSUBLISTENERCOMMAND")
	}()

	//Подключение к локальному FTP серверу
	os.Setenv("GO_PHFTP_LOCALFTP_HOST", "34.56.232.5")
	os.Setenv("GO_PHFTP_LOCALFTP_USERNAME", "local_user_name")
	defer func() {
		os.Unsetenv("GO_PHFTP_LOCALFTP_HOST")
		os.Unsetenv("GO_PHFTP_LOCALFTP_USERNAME")
	}()

	//Подключение к FTP серверу агрегатору
	os.Setenv("GO_PHFTP_MAINFTP_HOST", "67.43.123.33")
	os.Setenv("GO_PHFTP_MAINFTP_USERNAME", "main_user_name")
	defer func() {
		os.Unsetenv("GO_PHFTP_MAINFTP_HOST")
		os.Unsetenv("GO_PHFTP_MAINFTP_USERNAME")
	}()

	conf, err = confighandler.New("placeholder_ftp", CONF_DIR_NAME)
	assert.NoError(t, err)

	confNats = conf.GetConfigNATS()
	assert.Equal(t, confNats.Prefix, "new_prefix")
	assert.Equal(t, confNats.Host, "localhost")
	assert.Equal(t, confNats.Port, 3344)
	assert.Equal(t, confNats.CacheTTL, 4800)
	assert.Equal(t, confNats.Subscriptions.ListenerCommand, "phftp.commands.test")

	confLocalFTP = conf.GetConfigLocalFTP()
	assert.Equal(t, confLocalFTP.Host, "34.56.232.5")
	assert.Equal(t, confLocalFTP.Username, "local_user_name")

	confMainFTP = conf.GetConfigMainFTP()
	assert.Equal(t, confMainFTP.Host, "67.43.123.33")
	assert.Equal(t, confMainFTP.Username, "main_user_name")

}

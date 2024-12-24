package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_ftp/internal/wrappers"
)

var (
	srcFtp, dstFtp *wrappers.WrapperSimplyNetworkClient

	err error
)

type ConfFtp struct {
	username string
	password string
	host     string
	port     int
}

func NewConfFtp(host string, port int, isLocalFtp bool) (*ConfFtp, error) {
	conf := &ConfFtp{}

	if err := godotenv.Load(".env"); err != nil {
		return conf, err
	}

	username := os.Getenv("GO_PHFTP_MAINFTP_USERNAME")
	passwd := os.Getenv("GO_PHFTP_MAINFTP_PASSWD")

	if isLocalFtp {
		username = os.Getenv("GO_PHFTP_LOCALFTP_USERNAME")
		passwd = os.Getenv("GO_PHFTP_LOCALFTP_PASSWD")
	}

	if host == "" || port == 0 {
		return conf, errors.New("invalide parameter 'host' or 'port'")
	}

	conf.host = host
	conf.port = port

	if username == "" || passwd == "" {
		return conf, errors.New("'username' or 'passwd' it should not be empty")
	}

	conf.password = passwd
	conf.username = username

	return conf, nil
}

func (conf *ConfFtp) GetHost() string {
	return conf.host
}

func (conf *ConfFtp) GetPort() int {
	return conf.port
}

func (conf *ConfFtp) GetUsername() string {
	return conf.username
}

func (conf *ConfFtp) GetPasswd() string {
	return conf.password
}

func (conf *ConfFtp) SetPort(v int) {
	conf.port = v
}

func (conf *ConfFtp) SetHost(v string) {
	conf.host = v
}

func (conf *ConfFtp) SetPasswd(v string) {
	conf.username = v
}

func (conf *ConfFtp) SetUsername(v string) {
	conf.password = v
}

func TestMain(m *testing.M) {
	confLocalFtp, err := NewConfFtp("127.0.0.1", 21, true)
	if err != nil {
		log.Fatalln(err)
	}

	confMainFtp, err := NewConfFtp("ftp-users.cloud.gcm", 21, false)
	if err != nil {
		log.Fatalln(err)
	}

	srcFtp, err = wrappers.NewWrapperSimpleNetworkClient(confLocalFtp)
	if err != nil {
		log.Fatalln("src ftp error:", err)
	}

	dstFtp, err = wrappers.NewWrapperSimpleNetworkClient(confMainFtp)
	if err != nil {
		log.Fatalln("dst ftp error:", err)
	}

	os.Exit(m.Run())
}

func TestRequestFTPServer(t *testing.T) {
	t.Run("Test 1. Соединения с ftp", func(t *testing.T) {
		assert.NoError(t, srcFtp.CheckConn())
		assert.NoError(t, dstFtp.CheckConn())
	})

	t.Run("Test 2. Скачивание файла с ftp", func(t *testing.T) {
		num, err := srcFtp.ReadFile(context.Background(), wrappers.WrapperReadWriteFileOptions{
			SrcFilePath: "/ftp/someuser/folder_one",
			SrcFileName: "test_pcap_file.pcap",
			DstFilePath: "../../tmp_files/",
			DstFileName: "test_pcap_file.pcap",
		})

		fmt.Println("111 ERROR:", err)

		assert.NoError(t, err)
		assert.Greater(t, num, 1)
	})

	t.Run("Test 3. Загрузка файла на ftp", func(t *testing.T) {
		_, err := os.Stat(path.Join("../../tmp_files/", "test_pcap_file.pcap"))
		assert.False(t, errors.Is(err, os.ErrNotExist))

		err = dstFtp.WriteFile(context.Background(), wrappers.WrapperReadWriteFileOptions{
			SrcFilePath: "../../tmp_files/",
			SrcFileName: "test_pcap_file.pcap",
			DstFilePath: "/net_traff_txt",
			DstFileName: "_test_pcap_file.pcap",
		})

		fmt.Println("222 ERROR:", err)

		assert.NoError(t, err)

		//удаление файла из временной директории
		//err = os.Remove(path.Join("../../tmp_files/", "book.pdf"))
		//assert.NoError(t, err)
	})
}

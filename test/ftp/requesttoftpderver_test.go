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

func NewConfFtp(username, host string, port int) (*ConfFtp, error) {
	conf := &ConfFtp{}

	if err := godotenv.Load(".env"); err != nil {
		return conf, err
	}

	if username == "" || host == "" || port == 0 {
		return conf, errors.New("invalide parameter 'username', 'host' or 'port'")
	}

	conf.host = host
	conf.username = username
	conf.port = port

	passwd := os.Getenv("GO_PHFTP_LOCALFTP_PASSWD")
	if passwd == "" {
		return conf, errors.New("'passwd' it should not be empty")
	}

	conf.password = passwd

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
	conf, err := NewConfFtp("someuser", "127.0.0.1", 21)
	if err != nil {
		log.Fatalln(err)
	}

	srcFtp, err = wrappers.NewWrapperSimpleNetworkClient(conf)
	if err != nil {
		log.Fatalln("src ftp error:", err)
	}

	dstFtp, err = wrappers.NewWrapperSimpleNetworkClient(conf)
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
			SrcFileName: "book.pdf1",
			DstFilePath: "../../tmp_files/",
			DstFileName: "book.pdf",
		})

		fmt.Println("ERROR:", err)

		assert.NoError(t, err)
		assert.Greater(t, num, 1)
	})

	t.Run("Test 3. Загрузка файла на ftp", func(t *testing.T) {
		_, err := os.Stat(path.Join("../../tmp_files/", "book.pdf"))
		assert.False(t, errors.Is(err, os.ErrNotExist))

		err = dstFtp.WriteFile(context.Background(), wrappers.WrapperReadWriteFileOptions{
			SrcFilePath: "../../tmp_files/",
			SrcFileName: "book.pdf",
			DstFilePath: "/ftp/someuser/folder_two",
			DstFileName: "book_new.pdf",
		})

		fmt.Println("ERROR:", err)

		assert.NoError(t, err)

		//удаление файла из временной директории
		err = os.Remove(path.Join("../../tmp_files/", "book.pdf"))
		assert.NoError(t, err)
	})
}

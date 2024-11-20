package wrappers

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/jlaffaye/ftp"
)

//********** WrapperSimplyNetworkClient ************

func (client *WrapperSimplyNetworkClient) setHost(v string) {
	client.host = v
}

func (client *WrapperSimplyNetworkClient) setPort(v int) {
	client.port = v
}

func (client *WrapperSimplyNetworkClient) setUsername(v string) {
	client.username = v
}

func (client *WrapperSimplyNetworkClient) setPasswd(v string) {
	client.passwd = v
}

// CheckConn проверка наличие сетевого доступа
func (client *WrapperSimplyNetworkClient) CheckConn() error {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", client.host, client.port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer func(c *ftp.ServerConn, err error) {
		if errQuite := c.Quit(); errQuite != nil {
			errors.Join(err, errQuite)
		}
	}(c, err)

	err = c.Login(client.username, client.passwd)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile чтение файла
func (client *WrapperSimplyNetworkClient) ReadFile(filePath, fileName string) ([]byte, int, error) {
	var (
		byteFile []byte = make([]byte, 0)
		c        *ftp.ServerConn

		err error
	)

	c, err = ftp.Dial(fmt.Sprintf("%s:%d", client.host, client.port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return byteFile, 0, err
	}
	defer func(c *ftp.ServerConn, err error) {
		if errQuite := c.Quit(); errQuite != nil {
			errors.Join(err, errQuite)
		}
	}(c, err)

	err = c.Login(client.username, client.passwd)
	if err != nil {
		return byteFile, 0, err
	}

	r, err := c.Retr(path.Join(filePath, fileName))
	if err != nil {
		return byteFile, 0, err
	}
	defer func(r *ftp.Response, err error) {
		if errClose := r.Close(); errClose != nil {
			errors.Join(err, errClose)
		}
	}(r, err)

	if _, err = r.Read(byteFile); err != nil {
		return byteFile, 0, err
	}

	return byteFile, 0, nil
}

// WriteFile запись файла
func (client *WrapperSimplyNetworkClient) WriteFile() {

}

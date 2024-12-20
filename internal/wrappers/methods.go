package wrappers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
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
func (client *WrapperSimplyNetworkClient) ReadFile(ctx context.Context, opts WrapperReadWriteFileOptions) (int, error) {
	var size int

	c, err := ftp.Dial(
		fmt.Sprintf("%s:%d", client.host, client.port),
		ftp.DialWithContext(ctx),
		ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return size, err
	}
	defer func(c *ftp.ServerConn, err error) {
		if errQuite := c.Quit(); errQuite != nil {
			errors.Join(err, errQuite)
		}
	}(c, err)

	err = c.Login(client.username, client.passwd)
	if err != nil {
		return size, err
	}

	r, err := c.Retr(path.Join(opts.SrcFilePath, opts.SrcFileName))
	if err != nil {
		return size, err
	}
	defer func(r *ftp.Response, err error) {
		if errClose := r.Close(); errClose != nil {
			errors.Join(err, errClose)
		}
	}(r, err)

	f, err := os.OpenFile(path.Join(opts.DstFilePath, opts.DstFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return size, err
	}
	defer func(f *os.File, err error) {
		if errClose := f.Close(); errClose != nil {
			errors.Join(err, errClose)
		}
	}(f, err)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		num, err := f.Write(b)
		if err != nil {
			return size, err
		}

		size += num
	}

	return size, nil
}

// WriteFile запись файла
func (client *WrapperSimplyNetworkClient) WriteFile(ctx context.Context, opts WrapperReadWriteFileOptions) error {
	filePath := path.Join(opts.SrcFilePath, opts.SrcFileName)
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File, err error) {
		if errClose := f.Close(); errClose != nil {
			errors.Join(err, errClose)
		}
	}(f, err)

	c, err := ftp.Dial(
		fmt.Sprintf("%s:%d", client.host, client.port),
		ftp.DialWithContext(ctx),
		ftp.DialWithTimeout(5*time.Second))
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

	err = c.Stor(path.Join(opts.DstFilePath, opts.DstFileName), f)
	if err != nil {
		return err
	}
	defer func(c *ftp.ServerConn, err error) {
		if errQuite := c.Quit(); errQuite != nil {
			errors.Join(err, errQuite)
		}
	}(c, err)

	return nil
}

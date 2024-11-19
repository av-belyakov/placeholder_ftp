package wrappers

//********** WrapperFtpClient ************

func (client *WrapperFtpClient) setHost(v string) {
	client.host = v
}

func (client *WrapperFtpClient) setPort(v int) {
	client.port = v
}

func (client *WrapperFtpClient) setUsername(v string) {
	client.username = v
}

func (client *WrapperFtpClient) setPasswd(v string) {
	client.passwd = v
}

/*
	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

здесь надо описать методы чтения файла с FTP серврера и
записи файла на FTP сервер
*/
func (client *WrapperFtpClient) ReadFile() {

}

func (client *WrapperFtpClient) WriteFile() {

}

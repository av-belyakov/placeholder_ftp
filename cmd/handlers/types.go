package handlers

// RequestCopyFileFromFtpServer структура запроса для обработки файлов на FTP сервере
type RequestCopyFileFromFtpServer struct {
	TaskId     string                         `json:"task_id"`
	Service    string                         `json:"service"`
	Command    string                         `json:"command"`
	Parameters ParameterCopyFileFromFtpServer `json:"parameters"`
}

// ParameterCopyFileFromFtpServer подробные параметры
type ParameterCopyFileFromFtpServer struct {
	PathLocalFtp string   `json:"path_local_ftp"`
	PathMainFtp  string   `json:"path_main_ftp"`
	Files        []string `json:"files"`
}

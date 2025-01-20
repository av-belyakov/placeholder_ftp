package supportingfunctions

import (
	"net/url"
	"path"
	"strings"
)

type LinkParseResult struct {
	Scheme   string //тип протокола http, ftp
	Host     string
	Path     string
	FileName string
}

// LinkParse разбирает ссылку на фрагменты
func LinkParse(urlStr string) (LinkParseResult, error) {
	result := LinkParseResult{}

	baseUrl, err := url.Parse(urlStr)
	if err != nil {
		return result, err
	}

	tmp := strings.Split(baseUrl.Path, "/")

	result.Scheme = baseUrl.Scheme
	result.Host = baseUrl.Host
	result.Path = path.Join(tmp[:(len(tmp) - 1)]...)
	result.FileName = tmp[len(tmp)-1]

	return result, nil
}

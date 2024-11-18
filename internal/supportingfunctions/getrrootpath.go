package supportingfunctions

import (
	"os"
	"strings"
)

func GetRootPath(rootDir string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tmp := strings.Split(currentDir, "/")

	if tmp[len(tmp)-1] == rootDir {
		return currentDir, nil
	}

	var path string = ""
	for _, v := range tmp {
		path += v + "/"

		if v == rootDir {
			return path, nil
		}
	}

	return path, nil
}

package supportingfunctions

import (
	"os"
	"path"
)

func CreateDirectory(rootDir, dirName string) error {
	newPath, err := GetRootPath(rootDir)
	if err != nil {
		return err
	}

	pd := path.Join(newPath, dirName)

	if _, err := os.ReadDir(pd); err != nil {
		if err := os.Mkdir(pd, 0777); err != nil {
			return err
		}
	}

	return nil
}

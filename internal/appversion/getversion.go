package appversion

import "os"

// GetAppVersion версия приложения
func GetAppVersion() (string, error) {
	b, err := os.ReadFile("version")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

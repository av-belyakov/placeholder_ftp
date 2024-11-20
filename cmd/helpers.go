package main

import (
	"os"
	"path"

	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
	"github.com/av-belyakov/simplelogger"
)

func getLoggerSettings(cls []confighandler.LoggerOption) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			MsgTypeName:     v.MsgTypeName,
			WritingToFile:   v.WritingFile,
			PathDirectory:   v.PathDirectory,
			WritingToStdout: v.WritingStdout,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}

func createTmpDirectory(rootDir, dirName string) error {
	pd := path.Join(rootDir, dirName)

	if _, err := os.ReadDir(pd); err != nil {
		if err := os.Mkdir(pd, 0777); err != nil {
			return err
		}
	}

	return nil
}

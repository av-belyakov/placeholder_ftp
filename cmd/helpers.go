package main

import (
	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
	"github.com/av-belyakov/simplelogger"
)

func getLoggerSettings(cls []confighandler.LoggerOption) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			WritingToStdout: v.WritingStdout,
			WritingToFile:   v.WritingFile,
			WritingToDB:     v.WritingDB,
			PathDirectory:   v.PathDirectory,
			MsgTypeName:     v.MsgTypeName,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}

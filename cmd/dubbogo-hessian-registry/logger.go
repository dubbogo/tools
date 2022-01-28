package main

import (
	"log"
)

type logType int

const (
	infoLog logType = iota
	errorLog
)

// showLog 输出日志
func showLog(t logType, format string, v ...interface{}) {
	format = format + "\n"
	if onlyError == nil || !(*onlyError) || t == errorLog {
		log.Printf(format, v...)
		return
	}
}

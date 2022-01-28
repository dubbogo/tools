package generator

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
	if t == errorLog {
		log.Fatalf(format, v...)
		return
	}
	log.Printf(format, v...)
}

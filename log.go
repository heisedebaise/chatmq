package chatmq

import (
	"log"
)

const (
	//LogLevelDebug log level debug.
	LogLevelDebug LogLevel = 0
	//LogLevelInfo log level info.
	LogLevelInfo LogLevel = 1
	//LogLevelWarn log level warn.
	LogLevelWarn LogLevel = 2
	//LogLevelOff log level off.
	LogLevelOff LogLevel = 3
)

//LogLevel log level
type LogLevel int

var logLevels = []string{"DEBUG ", "INFO  ", "WARN  "}
var logLevel = LogLevelWarn

//SetLogLevel set log level.
func SetLogLevel(level LogLevel) {
	logLevel = level
}

func logf(level LogLevel, format string, v ...interface{}) {
	if level >= logLevel {
		log.Printf("chatmq: "+logLevels[level]+format+"\n", v...)
	}
}

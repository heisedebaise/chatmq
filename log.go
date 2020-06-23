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

func debug(format string, v ...interface{}) {
	if logLevel <= LogLevelDebug {
		logf("DEBUG ", format, v...)
	}
}

func info(format string, v ...interface{}) {
	if logLevel <= LogLevelInfo {
		logf("INFO  ", format, v...)
	}
}

func warn(format string, v ...interface{}) {
	if logLevel <= LogLevelWarn {
		logf("WARN  ", format, v...)
	}
}

func logf(level, format string, v ...interface{}) {
	log.Printf("chatmq: "+level+format+"\n", v...)
}

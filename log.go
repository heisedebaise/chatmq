package chatmq

import (
	"log"
)

func logf(format string, v ...interface{}) {
	log.Printf("chatmq: "+format, v...)
}

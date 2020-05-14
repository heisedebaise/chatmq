package chatmq

import (
	"time"
)

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
			queueOverdue()
			receiveOverdue()
			ping()
		}
	}()
}

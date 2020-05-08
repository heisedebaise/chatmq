package chatmq

import "time"

var overdue = time.Minute

//Overdue set overdue duration.
func Overdue(duration time.Duration) {
	overdue = duration
}

func clean() {
	go func() {
		for {
			time.Sleep(time.Second)
			t := time.Now().Add(-overdue)
			mq.Range(func(key, value interface{}) bool {
				if q, ok := value.(*queue); ok && q.time.Before(t) {
					mq.Delete(key)
				}

				return true
			})
		}
	}()
}

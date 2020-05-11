package chatmq

import (
	"crypto/md5"
	"sync"
	"time"
)

type queue struct {
	lock chan bool
	data [][]byte
	time time.Time
}

var mq sync.Map
var queueOverdueDuration = time.Minute

func skey(skey string) [16]byte {
	return key([]byte(skey))
}

func key(key []byte) [16]byte {
	return md5.Sum(key)
}

func put(key [16]byte, data []byte) {
	if v, ok := mq.Load(key); ok {
		putv(key, v, data)

		return
	}

	q := newQueue()
	if v, ok := mq.LoadOrStore(key, q); ok {
		putv(key, v, data)
	} else {
		putq(q, data)
	}
}

func putv(k, v interface{}, data []byte) {
	if q, ok := v.(*queue); ok {
		putq(q, data)
	} else {
		q := newQueue()
		mq.Store(k, q)
		putq(q, data)
	}
}

func newQueue() *queue {
	return &queue{lock: make(chan bool, 1), data: make([][]byte, 0), time: time.Now()}
}

func putq(q *queue, data []byte) {
	q.lock <- true
	q.data = append(q.data, data)
	<-q.lock
}

func get(key [16]byte) (data []byte) {
	if v, ok := mq.Load(key); ok {
		if q, ok := v.(*queue); ok {
			q.lock <- true
			if len(q.data) > 0 {
				data = q.data[0]
				q.data = q.data[1:]
			}
			q.time = time.Now()
			<-q.lock
		} else {
			mq.Delete(key)
		}
	}

	return
}

//QueueOverdue set queue overdue duration.
func QueueOverdue(duration time.Duration) {
	queueOverdueDuration = duration
}

func queueOverdue() {
	t := time.Now().Add(-queueOverdueDuration)
	mq.Range(func(key, value interface{}) bool {
		if q, ok := value.(*queue); ok && q.time.Before(t) {
			mq.Delete(key)
		}

		return true
	})
}

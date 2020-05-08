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

func skey(skey string) [16]byte {
	return key([]byte(skey))
}

func key(key []byte) [16]byte {
	return md5.Sum(key)
}

func put(key [16]byte, data []byte) {
	if v, ok := mq.Load(key); ok {
		putv(v, data)

		return
	}

	q := &queue{lock: make(chan bool, 1), data: make([][]byte, 0)}
	if v, ok := mq.LoadOrStore(key, q); ok {
		putv(v, data)
	} else {
		putq(q, data)
	}
}

func putv(v interface{}, data []byte) {
	if q, ok := v.(*queue); ok {
		putq(q, data)
	}
}

func putq(q *queue, data []byte) {
	q.lock <- true
	q.data = append(q.data, data)
	q.time = time.Now()
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
		}
	}

	return
}

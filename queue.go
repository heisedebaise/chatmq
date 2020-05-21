package chatmq

import (
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

//Put put.
func Put(key string, e interface{}) {
	data, err := encode(e)
	if err != nil {
		logf(LogLevelWarn, "encode %v fail %v", e, err)

		return
	}

	k := skey(key)
	put(k, data)
	nodes.Range(func(_, value interface{}) bool {
		if node, ok := value.(*node); ok {
			node.send(methodPut, k, data)
		}

		return true
	})
}

func put(key [16]byte, data []byte) {
	logf(LogLevelDebug, "put %v %d", key, len(data))
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

//Get get.
func Get(key string, e interface{}) bool {
	if data, ok := get(skey(key)); ok {
		err := Decode(data, e)
		if err == nil {
			return true
		}

		logf(LogLevelWarn, "decode fail %v", err)
	}

	return false
}

func get(key [16]byte) (data []byte, ok bool) {
	if v, o := mq.Load(key); o {
		if q, o := v.(*queue); o {
			q.lock <- true
			if len(q.data) > 0 {
				ok = true
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

//Overdue set queue overdue duration.
func Overdue(duration time.Duration) {
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

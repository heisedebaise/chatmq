package chatmq

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMQ(t *testing.T) {
	key := skey("chatmq")
	if data, ok := get(key); ok || len(data) > 0 {
		t.Errorf("not empty %v\n", data)
	}
	put(key, []byte("chat mq"))
	if data, ok := get(key); !ok || !bytes.Equal(data, []byte("chat mq")) {
		t.Errorf("get fail %s \n", string(data))
	}
	if data, ok := get(key); ok || len(data) > 0 {
		t.Errorf("not empty %v\n", data)
	}

	wg := sync.WaitGroup{}
	size := 64
	wg.Add(size << 1)
	for i := 0; i < size; i++ {
		ii := i
		go func() {
			put(key, []byte("data "+strconv.Itoa(ii)))
			wg.Done()
		}()
	}

	ns := make([]int, size)
	for i := 0; i < size; i++ {
		go func() {
			for {
				if data, ok := get(key); ok && len(data) > 5 {
					if n, err := strconv.Atoi(string(data[5:])); err == nil {
						ns[n]++
					}
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < size; i++ {
		if ns[i] != 1 {
			t.Errorf("fail %d %d!=1\n", i, ns[i])
		}
	}

	queueOverdueDuration = 2 * time.Second
	if _, ok := mq.Load(key); !ok {
		t.Errorf("key not exists\n")
	}

	time.Sleep(3 * time.Second)
	if _, ok := mq.Load(key); ok {
		t.Errorf("key exists\n")
	}
}

func TestPutGet(t *testing.T) {
	type data struct {
		Key   string
		Value string
	}

	d := data{}
	if Get("key", &d) || d.Key != "" || d.Value != "" {
		t.Errorf("get not empty data %v\n", d)
	}

	Put("key", data{"key", "value"})
	if !Get("key", &d) {
		t.Errorf("get no data\n")
	}
	if d.Key != "key" || d.Value != "value" {
		t.Errorf("get illegal data %v\n", d)
	}
}

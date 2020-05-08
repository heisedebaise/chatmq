package chatmq

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
)

func TestMQ(t *testing.T) {
	key := key([]byte("chatmq"))
	if data := get(key); len(data) > 0 {
		t.Errorf("not empty %v\n", data)
	}
	put(key, []byte("chat mq"))
	if data := get(key); !bytes.Equal(data, []byte("chat mq")) {
		t.Errorf("get fail %s \n", string(data))
	}
	if data := get(key); len(data) > 0 {
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
				data := get(key)
				if len(data) > 5 {
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
}

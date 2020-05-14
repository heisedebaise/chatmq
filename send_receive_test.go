package chatmq

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestSendReceive(t *testing.T) {
	go Cluster(":9371", "secret key", []string{":9371"})
	time.Sleep(time.Second)

	nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*node); ok {
			node.lock <- true
			node.state = 1
			<-node.lock
		}

		return true
	})

	key := []byte("key")
	testSendReceive(key, bufferSize>>1, t)
	testSendReceive(key, bufferSize, t)
	testSendReceive(key, bufferSize<<1, t)
}

func testSendReceive(key []byte, size int, t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, 0)
	for i := 0; i < size; i++ {
		data = append(data, byte(rand.Intn(0xff)))
	}

	Put(key, data)
	time.Sleep(time.Second)
	for i := 0; i < 2; i++ {
		if d, ok := Get(key); !ok || !bytes.Equal(d, data) {
			t.Errorf("illegal data %d %t %d %d\n", i, ok, len(data), len(d))
		}
	}
}
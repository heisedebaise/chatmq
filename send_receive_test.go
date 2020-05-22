package chatmq

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestSendReceive(t *testing.T) {
	testClusterUp()
	key := skey("key")
	testSendReceive(key, bufferSize>>1, time.Second, t)
	testSendReceive(key, bufferSize, time.Second<<1, t)
	testSendReceive(key, bufferSize<<1, time.Second<<2, t)
}

func testSendReceive(key [16]byte, size int, sleep time.Duration, t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, 0)
	for i := 0; i < size; i++ {
		data = append(data, byte(rand.Intn(0xff)))
	}

	testNoSelfListen()
	sends(methodPut, key, data)
	time.Sleep(sleep)
	if d, ok := get(key); !ok || !bytes.Equal(d, data) {
		t.Errorf("illegal data %t %d %d\n", ok, len(data), len(d))
	}
}

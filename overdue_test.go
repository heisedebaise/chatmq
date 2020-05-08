package chatmq

import (
	"testing"
	"time"
)

func TestOverdue(t *testing.T) {
	Overdue(3 * time.Second)

	key := skey("chatmq")
	if _, ok := mq.Load(key); ok {
		t.Errorf("key exists\n")
	}

	put(key, []byte("chat mq"))
	if _, ok := mq.Load(key); !ok {
		t.Errorf("key not exists\n")
	}

	get(key)
	if _, ok := mq.Load(key); !ok {
		t.Errorf("key not exists\n")
	}

	time.Sleep(5 * time.Second)
	if _, ok := mq.Load(key); ok {
		t.Errorf("key exists\n")
	}
}

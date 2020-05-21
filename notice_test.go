package chatmq

import (
	"testing"
	"time"
)

func TestNotice(t *testing.T) {
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

	var data string
	SetNotice("notice key", func(b []byte) {
		data = string(b)
	})

	SendNotice("notice key", []byte("notice data"))
	time.Sleep(time.Second)
	if data != "notice data" {
		t.Errorf("illegal data %s\n", data)
	}

	SendNotice("notice-key", []byte("notice-data"))
	time.Sleep(time.Second)
	if data != "notice data" {
		t.Errorf("illegal data %s\n", data)
	}
}

package chatmq

import (
	"testing"
	"time"
)

func TestNotice(t *testing.T) {
	testClusterUp()
	testNoSelfListen()
	var data string
	SetNotice("notice key", func(b []byte) {
		if err := Decode(b, &data); err != nil {
			t.Errorf("decode notice data fail %v\n", err)
		}
	})

	testNoSelfListen()
	SendNotice("notice key", "notice data")
	time.Sleep(time.Second)
	if data != "notice data" {
		t.Errorf("illegal data %s\n", data)
	}

	testNoSelfListen()
	SendNotice("notice-key", "notice-data")
	time.Sleep(time.Second)
	if data != "notice data" {
		t.Errorf("illegal data %s\n", data)
	}

	testNoSelfListen()
	DeleteNotice("notice key")
	SendNotice("notice key", "notice-data")
	time.Sleep(time.Second)
	if data != "notice data" {
		t.Errorf("illegal data %s\n", data)
	}

	testNoSelfListen()
	var n int
	SetNotice("nkey", func(b []byte) {
		if err := Decode(b, &n); err != nil {
			t.Errorf("decode notice data fail %v\n", err)
		}
	})
	SendNotice("nkey", 1)
	time.Sleep(time.Second)
	if n != 1 {
		t.Errorf("illegal data %d\n", n)
	}
}

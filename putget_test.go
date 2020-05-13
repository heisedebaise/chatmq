package chatmq

import (
	"bytes"
	"testing"
)

func TestPutGet(t *testing.T) {
	skey := "chatmq key"
	bkey := []byte(skey)
	if ok, data := Get(bkey); ok || len(data) != 0 {
		t.Errorf("get not empty data\n")
	}
	if ok, data := GetString(skey); ok || data != "" {
		t.Errorf("get not empty data\n")
	}

	Put(bkey, []byte("byte data"))
	PutString(skey, "string data")
	if ok, data := GetString(skey); !ok || data != "byte data" {
		t.Errorf("get data not equals\n")
	}
	if ok, data := Get(bkey); !ok || !bytes.Equal(data, []byte("string data")) {
		t.Errorf("get data not equals\n")
	}
}

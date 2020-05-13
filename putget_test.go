package chatmq

import (
	"bytes"
	"testing"
)

func TestPutGet(t *testing.T) {
	skey := "chatmq key"
	bkey := []byte(skey)
	if len(Get(bkey)) != 0 || GetString(skey) != "" {
		t.Errorf("get not empty data\n")
	}

	Put(bkey, []byte("byte data"))
	PutString(skey, "string data")
	if GetString(skey) != "byte data" || !bytes.Equal(Get(bkey), []byte("string data")) {
		t.Errorf("get data not equals\n")
	}
}

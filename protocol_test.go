package chatmq

import (
	"bytes"
	"testing"
)

func TestPack(t *testing.T) {
	id := newID()
	data := []byte("chatmq data")
	key := skey("chatmq key")
	p := pack(id, data, 1, 2, 3, key)
	if !bytes.Equal(getID(p), id) {
		t.Errorf("id not equals\n")
	}
	if getSize(p) != 1 {
		t.Errorf("size not equals\n")
	}
	if getIndex(p) != 2 {
		t.Errorf("index not equals\n")
	}
	if getMethod(p) != 3 {
		t.Errorf("method not equals\n")
	}
	if getKey(p) != key {
		t.Errorf("key not equals\n")
	}
	if !bytes.Equal(getData(p), data) {
		t.Errorf("data not equals\n")
	}
}

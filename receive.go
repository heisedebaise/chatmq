package chatmq

import (
	"bytes"
	"net"
	"time"
)

var receiveLock = make(chan bool, 1)
var receiveDatas = make(map[string]*receiveData)

type receiveData struct {
	lock  chan bool
	data  [][]byte
	index int
	time  time.Time
}

func (rd *receiveData) put(index uint32, b []byte) (done bool) {
	rd.lock <- true
	rd.data[index] = b
	rd.index++
	done = rd.index == len(rd.data)
	<-rd.lock

	return
}

func (rd *receiveData) get() []byte {
	if rd.index == 1 {
		return rd.data[0]
	}

	var buffer bytes.Buffer
	for _, data := range rd.data {
		buffer.Write(data)
	}

	return buffer.Bytes()
}

func receive(conn *net.UDPConn, addr *net.UDPAddr, b []byte) {
	decrypt, err := decrypt(b)
	if err != nil {
		return
	}

	length := len(decrypt)
	if length < minLength {
		return
	}

	id := getID(decrypt)
	sid := string(id)
	size := getSize(decrypt)
	rd := getReceiveData(sid, size)
	done := rd.put(getIndex(decrypt), decrypt[keyEnd:])
	if !done {
		return
	}

	switch getMethod(decrypt) {
	case 0:
		send(conn, addr, id, listenID, 0, emptyKey)
	case 1:
		put(getKey(decrypt), rd.get())
	}

	if size > 1 {
		receiveLock <- true
		delete(receiveDatas, sid)
		<-receiveLock
	}
}

func getReceiveData(id string, size uint32) *receiveData {
	if size <= 1 {
		return &receiveData{make(chan bool, 1), make([][]byte, 1), 0, time.Now()}
	}

	receiveLock <- true
	rd, ok := receiveDatas[id]
	if !ok {
		rd = &receiveData{make(chan bool, 1), make([][]byte, int(size)), 0, time.Now()}
		receiveDatas[id] = rd
	}
	<-receiveLock

	return rd
}

func receiveOverdue() {
	if len(receiveDatas) == 0 {
		return
	}

	receiveLock <- true
	defer func() { <-receiveLock }()

	t := time.Now().Add(-time.Minute)
	ids := make([]string, 0)
	for id := range receiveDatas {
		if rd := receiveDatas[id]; rd.time.Before(t) {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return
	}

	for _, id := range ids {
		delete(receiveDatas, id)
	}
}

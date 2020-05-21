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
	logf(LogLevelDebug, "receive data %t %d/%d %d", done, index, len(rd.data), len(b))

	return
}

func (rd *receiveData) get() (data []byte) {
	if rd.index == 1 {
		data = rd.data[0]
	} else {
		var buffer bytes.Buffer
		for _, data := range rd.data {
			buffer.Write(data)
		}
		data = buffer.Bytes()
	}
	logf(LogLevelDebug, "get receive data %d %d", rd.index, len(rd.data[0]))

	return
}

func receive(conn *net.UDPConn, addr *net.UDPAddr, b []byte) {
	decrypt, err := decrypt(b)
	if err != nil {
		logf(LogLevelWarn, "decrypt fail %v", err)

		return
	}

	length := len(decrypt)
	if length < minLength {
		logf(LogLevelWarn, "illegal length %d<%d", length, minLength)

		return
	}

	id := getID(decrypt)
	sid := string(id)
	size := getSize(decrypt)
	rd := getReceiveData(sid, size)
	done := rd.put(getIndex(decrypt), getData(decrypt))
	if !done {
		return
	}

	switch getMethod(decrypt) {
	case methodPing:
		send(conn, addr, id, listenID, 0, emptyKey)
	case methodPut:
		put(getKey(decrypt), rd.get())
	case methodNotice:
		notice(getKey(decrypt), rd.get())
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

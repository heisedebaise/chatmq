package chatmq

import (
	"bytes"
	"encoding/binary"

	"github.com/google/uuid"
)

const (
	bufferSize  = 64 * 1024
	idEnd       = 36
	sizeEnd     = idEnd + 4
	indexEnd    = sizeEnd + 4
	methodEnd   = indexEnd + 1
	keyEnd      = methodEnd + 16
	minLength   = keyEnd
	dataMaxSize = bufferSize - cryptSize - minLength - 16
)

const (
	methodPing   method = 0
	methodPut    method = 1
	methodNotice method = 2
)

var emptyKey [16]byte

type method byte

func newID() []byte {
	return []byte(uuid.New().String())
}

func getID(b []byte) []byte {
	return b[:idEnd]
}

func getSize(b []byte) uint32 {
	return binary.BigEndian.Uint32(b[idEnd:sizeEnd])
}

func getIndex(b []byte) uint32 {
	return binary.BigEndian.Uint32(b[sizeEnd:indexEnd])
}

func getMethod(b []byte) method {
	return method(b[indexEnd])
}

func getKey(b []byte) (key [16]byte) {
	copy(key[:], b[methodEnd:keyEnd])

	return
}

func getData(b []byte) []byte {
	return b[keyEnd:]
}

func pack(id, data []byte, size, index uint32, m method, key [16]byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(id)
	buffer.Write(uint32byte(size))
	buffer.Write(uint32byte(index))
	buffer.WriteByte(byte(m))
	buffer.Write(key[:])
	buffer.Write(data)

	return buffer.Bytes()
}

func uint32byte(ui uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ui)

	return b
}

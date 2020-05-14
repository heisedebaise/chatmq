package chatmq

import (
	"net"
)

func send(conn *net.UDPConn, addr *net.UDPAddr, id, data []byte, method byte, key [16]byte) (int, error) {
	length := len(data)
	if length <= dataMaxSize {
		return sendTo(conn, addr, id, data, 0, 0, method, key)
	}

	size := length / dataMaxSize
	if size*dataMaxSize < length {
		size++
	}
	for i := 0; i < size; i++ {
		max := (i + 1) * dataMaxSize
		if max > length {
			max = length
		}
		n, err := sendTo(conn, addr, id, data[i*dataMaxSize:max], size, i, method, key)
		logf("send to udp %v %d %v\n", addr, length, err)
		if err != nil {
			return i*dataMaxSize + n, err
		}
	}

	return length, nil
}

func sendTo(conn *net.UDPConn, addr *net.UDPAddr, id, data []byte, size, index int, method byte, key [16]byte) (int, error) {
	p := pack(id, data, uint32(size), uint32(index), method, key)
	e, err := encrypt(p)
	if err != nil {
		return 0, err
	}

	return conn.WriteToUDP(e, addr)
}

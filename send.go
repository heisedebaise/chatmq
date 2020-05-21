package chatmq

import (
	"net"
)

func send(conn *net.UDPConn, addr *net.UDPAddr, id, data []byte, m method, key [16]byte) (int, error) {
	length := len(data)
	if length <= dataMaxSize {
		return sendTo(conn, addr, id, data, 0, 0, m, key)
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
		n, err := sendTo(conn, addr, id, data[i*dataMaxSize:max], size, i, m, key)
		if err != nil {
			return i*dataMaxSize + n, err
		}
	}

	return length, nil
}

func sendTo(conn *net.UDPConn, addr *net.UDPAddr, id, data []byte, size, index int, m method, key [16]byte) (int, error) {
	p := pack(id, data, uint32(size), uint32(index), m, key)
	e, err := encrypt(p)
	if err != nil {
		logf(LogLevelWarn, "encrypt fail %v", err)

		return 0, err
	}

	n, err := conn.WriteToUDP(e, addr)
	logf(LogLevelDebug, "send to udp %v %d/%d %d %v", addr, index, size, n, err)

	return n, err
}

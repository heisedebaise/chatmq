package chatmq

import (
	"net"
)

var buffer = make([]byte, 64*1024)
var listenID = newID()

func listen(host string) error {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	logf("listening on %s", host)
	for {
		n, a, err := conn.ReadFromUDP(buffer)
		logf("read from udp %v %d %v", a, n, err)
		if n <= 0 || err != nil {
			continue
		}

		data := make([]byte, n)
		copy(data, buffer)
		go receive(conn, a, data)
	}
}

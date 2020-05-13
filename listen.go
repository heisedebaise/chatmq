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

	logf("listening on %s\n", host)
	for {
		n, a, err := conn.ReadFromUDP(buffer)
		if n <= 0 || err != nil {
			logf("read from udp %v fail %d %v\n", a, n, err)

			continue
		}

		go receive(conn, a, buffer[:n])
	}
}

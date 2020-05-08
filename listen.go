package chatmq

import (
	"net"
)

var host = ":9371"
var buffer = make([]byte, 64*1024)

func listen() error {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	logf("chatmq listening on %s\n", host)
	for {
		n, a, err := conn.ReadFromUDP(buffer)
		if n <= 0 || err != nil {
			logf("read from udp %v fail %d %v\n", a, n, err)

			continue
		}

		go receive(a, buffer[:n])
	}
}

func receive(addr *net.UDPAddr, b []byte) {
}

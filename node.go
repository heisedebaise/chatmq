package chatmq

import (
	"bytes"
	"log"
	"net"
	"sync"
)

var nodes sync.Map

type node struct {
	addr  *net.UDPAddr
	conn  *net.UDPConn
	state int
}

func (n *node) check() {
	if n.state > 0 {
		return
	}

	id := newID()
	e, err := encrypt(pack(id, []byte{}, 0, 0, 0, emptyKey))
	if err != nil {
		return
	}

	n.conn.WriteToUDP(e, n.addr)
	buffer := make([]byte, bufferSize)
	c, _, err := n.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	d, err := decrypt(buffer[:c])
	if err != nil || len(d) < minLength || getMethod(d) != 0 || !bytes.Equal(getID(d), id) {
		return
	}

	if bytes.Equal(getData(d), listenID) {
		n.state = 2
	} else {
		n.state = 1
	}
}

func (n *node) send(key [16]byte, data []byte) {
	if n.state != 1 {
		return
	}

	if _, err := send(n.conn, n.addr, newID(), data, 1, key); err != nil {
		n.state = 0
		log.Printf("send udp to %v fail %v\n", n.addr, err)
	}
}

func setNodes(hosts []string) {
	m := make(map[string]bool)
	for _, host := range hosts {
		addr, err := net.ResolveUDPAddr("udp", host)
		if err != nil {
			logf("resolve udp addr %s fail %v\n", host, err)

			continue
		}

		conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
		if err != nil {
			logf("dial udp %s fail %v\n", host, err)

			continue
		}

		node := &node{addr, conn, 0}
		nodes.Store(host, node)
		m[host] = true
	}

	nodes.Range(func(key, value interface{}) bool {
		if host, ok := value.(string); ok {
			if _, ok = m[host]; !ok {
				nodes.Delete(host)
			}
		}

		return true
	})
}

func nodeState() {
	nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*node); ok {
			go node.check()
		} else {
			nodes.Delete(key.(string))
		}

		return true
	})
}

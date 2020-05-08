package chatmq

import (
	"net"
	"strings"
	"sync"
)

var nodes sync.Map

func setNodes(hosts string) {
	m := make(map[string]bool)
	for _, host := range strings.Split(hosts, ",") {
		addr, err := net.ResolveUDPAddr("udp", host)
		if err != nil {
			logf("resolve udp addr %s fail %v\n", host, err)

			continue
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			logf("dial udp %s fail %v\n", host, err)

			continue
		}

		nodes.Store(host, &node{addr: addr, conn: conn})
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

type node struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func (n *node) write(b []byte) (int, error) {
	return n.conn.WriteToUDP(b, n.addr)
}

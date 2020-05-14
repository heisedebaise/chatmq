package chatmq

import (
	"bytes"
	"net"
	"sync"
	"time"
)

var nodes sync.Map
var checkTimeout = 2 * time.Second

type node struct {
	addr  *net.UDPAddr
	conn  *net.UDPConn
	lock  chan bool
	state int // 0-unready;1-ready;2-self
	time  time.Time
}

func (n *node) ping() {
	if n.state == 2 || len(n.lock) > 0 || (n.state == 1 && n.time.Add(time.Minute).After(time.Now())) {
		return
	}

	id := newID()
	e, err := encrypt(pack(id, []byte{}, 0, 0, 0, emptyKey))
	if err != nil {
		return
	}

	n.lock <- true
	defer func() { <-n.lock }()

	n.conn.WriteToUDP(e, n.addr)
	buffer := make([]byte, minLength<<1)
	n.conn.SetReadDeadline(time.Now().Add(checkTimeout))
	c, addr, err := n.conn.ReadFromUDP(buffer)
	logf("ping udp %v %v %v %d %v\n", n.addr, checkTimeout, addr, c, err)
	if err != nil {
		n.state = 0

		return
	}

	d, err := decrypt(buffer[:c])
	if err != nil || len(d) < minLength || getMethod(d) != 0 || !bytes.Equal(getID(d), id) {
		n.state = 0

		return
	}

	if bytes.Equal(getData(d), listenID) {
		n.state = 2
	} else {
		n.state = 1
	}
	n.time = time.Now()
}

func (n *node) send(key [16]byte, data []byte) {
	if n.state != 1 {
		return
	}

	n.lock <- true
	defer func() { <-n.lock }()

	if _, err := send(n.conn, n.addr, newID(), data, 1, key); err != nil {
		n.state = 0
		logf("send udp to %v fail %v\n", n.addr, err)
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

		node := &node{addr, conn, make(chan bool, 1), 0, time.Now().Add(-time.Minute)}
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

func ping() {
	nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*node); ok {
			go node.ping()
		} else {
			nodes.Delete(key.(string))
		}

		return true
	})
}

//CheckTimeout check timeout.
func CheckTimeout(duration time.Duration) {
	checkTimeout = duration
}

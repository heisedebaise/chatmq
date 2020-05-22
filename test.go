package chatmq

import "time"

func testClusterUp() {
	go Cluster(":9371", "secret key", []string{":9371"})
	time.Sleep(time.Second)
}

func testNoSelfListen() {
	nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*node); ok {
			node.lock <- true
			node.state = 1
			<-node.lock
		}

		return true
	})
}

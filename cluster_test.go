package chatmq

import (
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	go Cluster(":9371", "secret key", []string{":9371"})
	time.Sleep(3 * time.Second)

	var size, self int
	nodes.Range(func(key, value interface{}) bool {
		size++
		if node, ok := value.(*node); ok && node.state == 2 {
			self++
		}

		return true
	})
	if size != 1 {
		t.Errorf("illegal size %d\n", size)
	}
	if self != 1 {
		t.Errorf("no self\n")
	}
}

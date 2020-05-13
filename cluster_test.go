package chatmq

import (
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	go Cluster(":9371", "secret key", []string{":9371", ":9372"})
	time.Sleep(5 * time.Second)

	var size, self, unready int
	nodes.Range(func(key, value interface{}) bool {
		size++
		if node, ok := value.(*node); ok {
			if node.state == 2 {
				self++
			} else if node.state == 0 || node.state == 3 {
				unready++
			}
		}

		return true
	})
	if size != 2 {
		t.Errorf("illegal size %d\n", size)
	}
	if self != 1 {
		t.Errorf("illegal self %d\n", self)
	}
	if unready != 1 {
		t.Errorf("illegal unready %d\n", unready)
	}
}

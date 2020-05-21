package chatmq

import (
	"sync"
)

var notices sync.Map

//Notice notice.
type Notice func(data []byte)

//SetNotice set notice.
func SetNotice(key string, notice Notice) {
	notices.Store(skey(key), notice)
}

//SendNotice send notice to other nodes.
func SendNotice(key string, data []byte) {
	sends(methodNotice, skey(key), data)
}

func notice(key [16]byte, data []byte) bool {
	if value, ok := notices.Load(key); ok {
		if notice, ok := value.(Notice); ok {
			notice(data)

			return true
		}
	}

	return false
}

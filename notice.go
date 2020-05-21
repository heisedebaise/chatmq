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

//DeleteNotice delete notice.
func DeleteNotice(key string) {
	notices.Delete(skey(key))
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

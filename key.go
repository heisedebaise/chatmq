package chatmq

import "crypto/md5"

func skey(key string) [16]byte {
	return md5.Sum([]byte(key))
}

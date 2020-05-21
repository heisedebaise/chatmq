package chatmq

import "crypto/md5"

func skey(skey string) [16]byte {
	return md5.Sum([]byte(skey))
}

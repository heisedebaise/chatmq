package chatmq

//Put put.
func Put(key, data []byte) {
	send2node(bkey(key), data)
}

//PutString put string.
func PutString(key, data string) {
	send2node(skey(key), []byte(data))
}

func send2node(key [16]byte, data []byte) {
	put(key, data)
	nodes.Range(func(k, v interface{}) bool {
		if node, ok := v.(*node); ok {
			node.send(key, data)
		}

		return true
	})
}

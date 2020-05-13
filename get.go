package chatmq

//Get get.
func Get(key []byte) ([]byte, bool) {
	return get(bkey(key))
}

//GetString get string.
func GetString(key string) (string, bool) {
	data, ok := get(skey(key))

	return string(data), ok
}

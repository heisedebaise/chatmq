package chatmq

//Get get.
func Get(key []byte) []byte {
	return get(bkey(key))
}

//GetString get string.
func GetString(key string) string {
	return string(get(skey(key)))
}

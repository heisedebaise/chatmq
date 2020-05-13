package chatmq

//Get get.
func Get(key []byte) (bool, []byte) {
	return get(bkey(key))
}

//GetString get string.
func GetString(key string) (bool, string) {
	ok, data := get(skey(key))

	return ok, string(data)
}

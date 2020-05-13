package chatmq

//Cluster cluster.
func Cluster(host, key string, nodes []string) error {
	cryptKey = key
	setNodes(nodes)

	return listen(host)
}

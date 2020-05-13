package chatmq

//Cluster cluster.
func Cluster(host, secret string, nodes []string) error {
	cryptKey = secret
	setNodes(nodes)

	return listen(host)
}

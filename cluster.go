package chatmq

//Cluster cluster.
func Cluster(host, secret string, nodes []string) error {
	cryptSecret = secret
	setNodes(nodes)

	return listen(host)
}

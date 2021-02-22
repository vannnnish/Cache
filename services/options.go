package services

type Options struct {
	Address string
	Port    int
	// ServerType 服务器类型
	ServerType string

	// VirtualNodeCount 一致性哈希虚拟节点个数
	VirtualNodeCount int

	// UpdateCircleDuration 更新一致性哈希的时间间隔
	UpdateCircleDuration int

	// Cluster 需要加入的集群
	Cluster []string
}

func DefaultOptions() Options {
	return Options{
		Address:              "127.0.0.1",
		Port:                 5837,
		ServerType:           "tcp",
		VirtualNodeCount:     1024,
		UpdateCircleDuration: 3,
		Cluster:              nil,
	}
}

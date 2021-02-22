package main

import (
	"cache/caches"
	"cache/services"
	"flag"
	"strings"
)

func main() {
	// 准备服务器的选项配置
	serverOptions := services.DefaultOptions()
	flag.StringVar(&serverOptions.Address, "address", serverOptions.Address, "The address used to listen, such as 127.0.0.1:5837")
	flag.IntVar(&serverOptions.Port, "prot", serverOptions.Port, "The port used to listen ,such as 5837")
	flag.StringVar(&serverOptions.ServerType, "serverType", serverOptions.ServerType, "The type of server (http ,tcp)")
	flag.IntVar(&serverOptions.VirtualNodeCount, "virtualNodeCount", serverOptions.VirtualNodeCount, "the number of virtual nodes in consistent hash")
	flag.IntVar(&serverOptions.UpdateCircleDuration, "updateCircleDuration", serverOptions.UpdateCircleDuration, "The duration between two circle updating operations. The unit is second.")

	cluster := flag.String("cluster", "", "The cluster of servers. One node in cluster will be ok")

	// 准备缓存配置选项
	options := caches.DefaultOptions()
	flag.IntVar(&options.MaxEntrySize, "maxEntrySize", options.MaxEntrySize, "The max memory size that entries can use . the unit is GB.")
	flag.IntVar(&options.MaxGcCount, "maxGcCount", options.MaxGcCount, "The")
	flag.Int64Var(&options.GcDuration, "gcDuration", options.GcDuration, "The duration between two gc tasks. The unit is Minute")

	// 获取持久化路径，和间隔时间
	flag.StringVar(&options.DumpFile, "dumpFile", options.DumpFile, "The file used to dump the cache")
	flag.Int64Var(&options.DumpDuration, "dumpDuration", options.DumpDuration, "The duration between two dump task")

	flag.IntVar(&options.MapSizeOfSegment, "mapSizeOfSegment", options.MapSizeOfSegment, "The map size of segment")
	flag.IntVar(&options.SegmentSize, "segmentSize", options.SegmentSize, "The number of segment in a cache. this value should be the pow of 2 for precision.")
	flag.IntVar(&options.CasSleepTime, "casSleepTime", options.CasSleepTime, "The time of sleep in one cas step. the unit is Microsecond")

	flag.Parse()

	// 从 flag 中解析出集群信息
	serverOptions.Cluster = nodesInCluster(*cluster)

	// 使用选项配置初始化缓存
	cache := caches.NewCacheWith(options)
	cache.AutoGc()
	cache.AutoDump()

	// 使用选项配置初始化服务器
	server, err := services.NewServer(cache, serverOptions)
	if err != nil {
		panic(err)
	}

	err = server.Run()
	if err != nil {
		panic(err)
	}
}

func nodesInCluster(cluster string) []string {
	if cluster == "" {
		return nil
	}
	return strings.Split(cluster, ",")
}

package main

import (
	"cache/caches"
	"cache/services"
	"flag"
	"log"
)

func main() {
	address := flag.String("address", ":5837", "The address used to listen, such as 127.0.0.1:5837")

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
	cache := caches.NewCacheWith(options)
	cache.AutoGc()
	cache.AutoDump()
	log.Printf("Kafo is runing on %s.", *address)
	err := services.NewHTTPServer(cache).Run(*address)
	if err != nil {
		panic(err)
	}
}

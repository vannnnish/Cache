package caches

type Options struct {
	// MaxEntrySize 键值对最大容量
	MaxEntrySize int
	// MaxGcCount 至每个segment 要清理的过期数据个数
	MaxGcCount int
	// GcDuration 多久执行一次 Gc 工作
	GcDuration int64

	// DumpFile 持久化文件路径
	DumpFile string
	// DumpDuration 持久化执行周期
	DumpDuration int64

	// MapSizeOfSegment segment 中map的初始化大小
	MapSizeOfSegment int

	// SegmentSize 缓存中有多少个segment
	SegmentSize int

	// CasSleepTime 每次CAS 自选需要等待时间 单位微妙
	CasSleepTime int
}

func DefaultOptions() Options {
	return Options{
		MaxEntrySize:     4,
		MaxGcCount:       10,
		GcDuration:       60,
		DumpFile:         "kafo.dump",
		DumpDuration:     30,
		MapSizeOfSegment: 256,
		SegmentSize:      1024,
		CasSleepTime:     1000,
	}
}

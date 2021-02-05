package caches

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type dump struct {
	SegmentSize int
	Segments    []*segment
	Options     *Options
}

func newEmptyDump() *dump {
	return &dump{}
}

func newDump(c *Cache) *dump {
	return &dump{
		SegmentSize: c.segmentSize,
		Options:     c.options,
		Segments:    c.segments,
	}
}

func nowSuffix() string {
	return "." + time.Now().Format("20060102150405")
}

func (d *dump) to(dumpFile string) error {
	newDumpFile := dumpFile + nowSuffix()
	file, err := os.OpenFile(newDumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	// 对文件进行gob编码
	err = gob.NewEncoder(file).Encode(d)
	if err != nil {
		file.Close()
		os.Remove(newDumpFile)
		return err
	}
	// 删除旧的持久化文件
	os.Remove(dumpFile)
	file.Close()
	return os.Rename(newDumpFile, dumpFile)
}

func (d *dump) from(dumpFile string) (*Cache, error) {
	file, err := os.Open(dumpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err = gob.NewDecoder(file).Decode(d); err != nil {
		return nil, err
	}
	// 恢复出 segment 之后需要为每个segment 的未导出字段进行初始化
	for _, segment := range d.Segments {
		segment.options = d.Options
		segment.lock = &sync.RWMutex{}
	}
	return &Cache{
		segmentSize: d.SegmentSize,
		options:     d.Options,
		segments:    d.Segments,
		dumping:     0,
	}, nil
}

package caches

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	// Value使用[]byte 是为了方便网络传输
	segmentSize int
	segments    []*segment
	//data        map[string]*value
	options *Options

	// dumping 表示当前缓存是不是处于持久化状态。 1表示处于持久化状态. 如果进入持久化状态，那么所有更新操作进入自旋状态，等待持久化完成
	dumping int32
	//status  *Status
	//lock    *sync.RWMutex
}

func NewCache() *Cache {
	return NewCacheWith(DefaultOptions())
}

func NewCacheWith(options Options) *Cache {
	if cache, ok := recoverFromDumpFile(options.DumpFile); ok {
		return cache
	}
	return &Cache{
		segmentSize: options.SegmentSize,
		segments:    newSegments(&options),
		options:     &options,
		dumping:     0,
	}
}

func recoverFromDumpFile(dumpFile string) (*Cache, bool) {
	cache, err := newEmptyDump().from(dumpFile)
	if err != nil {
		return nil, false
	}
	return cache, true
}

func newSegments(options *Options) []*segment {
	segments := make([]*segment, options.SegmentSize)
	for i := 0; i < options.SegmentSize; i++ {
		segments[i] = newSegment(options)
	}
	return segments
}

func index(key string) int {
	index := 0
	keyBytes := []byte(key)
	for _, b := range keyBytes {
		index = 31*index + int(b&0xff)
	}
	return index ^ (index >> 16)
}

func (c *Cache) segmentOf(key string) *segment {
	return c.segments[index(key)&(c.segmentSize-1)]
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.waitForDumping()
	return c.segmentOf(key).get(key)
}

func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, NeverDie)
}
func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.waitForDumping()
	return c.segmentOf(key).set(key, value, ttl)
}

func (c *Cache) Delete(key string) error {
	c.waitForDumping()
	c.segmentOf(key).delete(key)
	return nil
}

func (c *Cache) Status() Status {
	result := NewStatus()
	for _, segment := range c.segments {
		status := segment.status()
		result.Count += status.Count
		result.KeySize += status.KeySize
		result.ValueSize += status.ValueSize
	}
	return *result
}

//// 判断数据是否达到最大的容量
//func (c *Cache) checkEntrySize(newKey string, newValue []byte) bool {
//	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.options.MaxEntrySize*1024*1024
//}

// gc 会触发清理任务
func (c *Cache) gc() {
	c.waitForDumping()
	// 记录清理的个数
	wg := &sync.WaitGroup{}
	for _, segment := range c.segments {
		wg.Add(1)
		go func() {
			defer wg.Done()
			segment.gc()
		}()
	}
	wg.Wait()
}

func (c *Cache) AutoGc() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.GcDuration) * time.Minute)
		for {
			select {
			case <-ticker.C:
				c.gc()
			}
		}
	}()
}

func (c *Cache) dump() error {
	defer func() {
		fmt.Println("导出结束")
	}()
	atomic.StoreInt32(&c.dumping, 1)
	defer atomic.StoreInt32(&c.dumping, 0)
	return newDump(c).to(c.options.DumpFile)
}

func (c *Cache) AutoDump() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.DumpDuration) * time.Second)
		for {
			select {
			case <-ticker.C:
				c.dump()
			}
		}
	}()
}

func (c *Cache) waitForDumping() {
	for atomic.LoadInt32(&c.dumping) != 0 {
		time.Sleep(time.Duration(c.options.CasSleepTime) * time.Microsecond)
	}
}

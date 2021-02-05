package caches

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	concurrency = 100000
)

func testTask(task func(no int)) string {
	beginTime := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(no int) {
			defer wg.Done()
			task(no)
		}(i)
	}
	wg.Wait()
	return time.Now().Sub(beginTime).String()
}

func TestCacheSetGet(t *testing.T) {
	cache := NewCache()

	writeTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		cache.Set(data, []byte(data))
	})
	t.Logf("写入消耗时间为 %s", writeTime)
	time.Sleep(3 * time.Second)

	readTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		cache.Get(data)
	})
	t.Logf("读取消耗时间为: %s", readTime)
}

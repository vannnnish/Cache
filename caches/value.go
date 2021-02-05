package caches

import (
	"cache/helpers"
	"sync/atomic"
	"time"
)

const (
	NeverDie = 0
)

type value struct {
	Data  []byte
	Ttl   int64
	Ctime int64
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		Data:  helpers.Copy(data),
		Ttl:   ttl,
		Ctime: time.Now().Unix(),
	}
}

func (v *value) alive() bool {
	return v.Ttl == NeverDie || v.Ttl > time.Now().Unix()-v.Ctime
}

func (v *value) visit() []byte {
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}

package caches

import (
	"errors"
	"sync"
)

type segment struct {
	Data    map[string]*value
	Status  *Status
	options *Options
	lock    *sync.RWMutex
}

func newSegment(options *Options) *segment {
	return &segment{
		Data:    make(map[string]*value, options.MapSizeOfSegment),
		Status:  newStatus(),
		options: options,
		lock:    &sync.RWMutex{},
	}
}

func (s *segment) get(key string) ([]byte, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		return nil, false
	}
	if !value.alive() {
		s.lock.RUnlock()
		s.delete(key)
		s.lock.RLock()
		return nil, false
	}
	return value.visit(), true
}

func (s *segment) set(key string, value []byte, ttl int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
	}

	if !s.checkEntrySize(key, value) {
		if oldValue, ok := s.Data[key]; ok {
			s.Status.addEntry(key, oldValue.Data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}
	s.Status.addEntry(key, value)
	s.Data[key] = newValue(value, ttl)
	return nil
}

func (s *segment) delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
		delete(s.Data, key)
	}
}

func (s *segment) status() Status {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return *s.Status
}

func (s *segment) checkEntrySize(newKey string, newValue []byte) bool {
	return s.Status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= int64(s.options.MaxEntrySize*1024*2014/s.options.SegmentSize)
}

func (s *segment) gc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	count := 0
	for key, value := range s.Data {
		if !value.alive() {
			s.Status.subEntry(key, value.Data)
			delete(s.Data, key)
			count++
			if count >= s.options.MaxGcCount {
				break
			}
		}
	}
}

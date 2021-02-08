package caches

type Status struct {
	Count     int   `json:"count"`
	KeySize   int64 `json:"keySize"`
	ValueSize int64 `json:"valueSize"`
}

func NewStatus() *Status {
	return &Status{
		Count:     0,
		KeySize:   0,
		ValueSize: 0,
	}
}

func (s *Status) addEntry(key string, value []byte) {
	s.Count++
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(value))
}

func (s *Status) subEntry(key string, value []byte) {
	s.Count--
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}

func (s *Status) entrySize() int64 {
	return s.KeySize + s.ValueSize
}

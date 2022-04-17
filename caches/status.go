package caches


type Status struct {
	// numbers or key value pairs in cache
	Count int `json:"count"`
	// keySize record totol size of the key in cache
	KeySize int64 `json:"keySize"`
	// ValueSize record total size of the value in cache
	ValueSize int64 `json:"valueSize"`
}
func newStatus() *Status {
	return &Status{
		Count: 0,
		KeySize: 0,
		ValueSize: 0,
	}
}

func(s *Status) addEntry(key string, value []byte){
	s.Count ++
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(value))
}


func(s *Status) subEntry(key string, value []byte){
	s--
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}
// entrySize return sum of keySize and ValueSize
func (s *Status) entrySize() int64 {
    return s.KeySize + s.ValueSize
}
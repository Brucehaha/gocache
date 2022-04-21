package caches

import (
	"cache-server/helper"
	"sync/atomic"
	"time"
)

const (
	// NeverDie, TTL never expired if it is 0
	NeverDie = 0
)

// value is struct
type value struct {
	// data, store the real data
	data []byte
	// ttl is the time to expire
	ttl int64
	// ctime is the creation time of key value pairs
	ctime int64
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		data:  helper.Copy(data),
		ttl:   ttl,
		ctime: time.Now().Unix(),
	}
}

func (v *value) isAlive() bool {
	return v.ttl == NeverDie || time.Now().Unix()-v.ctime < v.ttl

}

//visit return the real data stored
func (v *value) visit() []byte {
	// use atomic instead of RLOCK, if RLOCK need the implement the write lock as well,
	// but lock will reduce the effiency, not use atomic.CompareAndSwapInt64 as
	// CAS will increase the cost the CPU
	atomic.SwapInt64(&v.ctime, time.Now().Unix())
	return v.data
}

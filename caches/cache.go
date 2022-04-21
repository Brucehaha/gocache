package caches

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	// store all value in map
	data map[string]*value
	// options not use pointer as it is just a config for data,
	// it will not be modified once passed in, and it is read only
	options Options
	// status wil be udpate frequently, use pointer type
	status *Status
	count  int64
	lock   *sync.RWMutex
}

func NewCache() *Cache {
	return NewCacheWith(DefaultOptions())
}

func NewCacheWith(options Options) *Cache {
	return &Cache{
		data:    make(map[string]*value, 256),
		options: options,
		status:  newStatus(),
		lock:    &sync.RWMutex{},
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if !value.isAlive() {
		// delete need WLOCK, need to release RLOCK
		c.lock.RUnlock()
		c.Delete(key)
		c.lock.RLock()
		return nil, false
	}
	return value.visit(), ok

}
func (c *Cache) Set(key string, value []byte) {

	c.SetWithTTL(key, value, NeverDie)
}

func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.checkEntrySize(key, value) {
		return errors.New("the entry size will exceed if you set this entry")
	}
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
	}
	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
	return nil
}

func (c *Cache) checkEntrySize(key string, value []byte) bool {
	return c.status.entrySize()+int64(len(key))+int64(len(value)) <= c.options.MaxEntrySize*1024*1024
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
		delete(c.data, key)
	}

}
func (c *Cache) Status() Status {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return *c.status
}
func (c *Cache) Count() int64 {
	c.lock.RLock()
	defer c.lock.RLock()
	return c.count
}

// gc will handle the cleaning of the expired data
func (c *Cache) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	count := 0
	for key, value := range c.data {
		if !value.isAlive() {
			c.status.subEntry(key, value.data)
			delete(c.data, key)
		}
		count++
		if count >= c.options.MaxGcCount {
			break
		}
	}
}

//AuoGc will start a timer for GC task
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

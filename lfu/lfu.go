// Based on https://github.com/dgrijalva/lfu-go/tree/master
package lfu

import (
	"container/list"
	"sync"
)

type Cache struct {
	values     map[string]*cacheEntry
	freqs      *list.List
	len        int
	lock       *sync.Mutex
	upperBound int
	lowerBound int
}

type cacheEntry struct {
	key      string
	value    interface{}
	freqNode *list.Element
}

type listEntry struct {
	entries map[*cacheEntry]byte
	freq    int
}

func New(upperBound int, lowerBound int) *Cache {
	c := new(Cache)
	c.values = make(map[string]*cacheEntry)
	c.freqs = list.New()
	c.lock = new(sync.Mutex)
	c.upperBound = upperBound
	c.lowerBound = lowerBound
	return c
}

func (c *Cache) Get(key string) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		e.value = value
		c.increment(e)
	} else {
		e := new(cacheEntry)
		e.key = key
		e.value = value
		c.values[key] = e
		c.increment(e)
		c.len++
		if c.len > c.upperBound {
			c.evict(c.len - c.lowerBound)
		}
	}
}

func (c *Cache) evict(count int) int {
	// No lock here so it can be called
	// from within the lock (during Set)
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry := range place.Value.(*listEntry).entries {
				if i < count {
					delete(c.values, entry.key)
					c.remEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *Cache) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next()
	}

	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		// create a new list entry
		li := new(listEntry)
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry]byte)
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry).entries[e] = 1
	if currentPlace != nil {
		// remove from current position
		c.remEntry(currentPlace, e)
	}
}

func (c *Cache) remEntry(place *list.Element, entry *cacheEntry) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}

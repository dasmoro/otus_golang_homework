package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}) bool
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[string]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity, items: map[string]*ListItem{}, queue: NewList()}
}

func (lru *lruCache) Get(key string) (interface{}, bool) {
	if val, ok := lru.items[key]; ok {
		v := val.Value.(cacheItem).value
		lru.queue.MoveToFront(val)
		return v, true
	}
	return nil, false
}

func (lru *lruCache) Set(key string, val interface{}) (have bool) {
	if v, ok := lru.items[key]; ok {
		v.Value = cacheItem{key, val}
		lru.queue.MoveToFront(v)
		have = true
	} else {
		if lru.capacity == lru.queue.Len() {
			last := lru.queue.Back()
			delete(lru.items, last.Value.(cacheItem).key)
			lru.queue.Remove(last)
		}
		lru.queue.PushFront(cacheItem{key, val})
		lru.items[key] = lru.queue.Front()
		have = false
	}
	return
}

func (lru *lruCache) Clear() {
	lru.queue = NewList()
	lru.items = map[string]*ListItem{}
}

type cacheItem struct {
	key   string
	value interface{}
}

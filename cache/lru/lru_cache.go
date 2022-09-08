package lru

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

type node[K constraints.Ordered, V any] struct {
	key    K
	val    V
	usedAt time.Time
}

// Cache is a LRU cache which is concurrent safe
type Cache[K constraints.Ordered, V any] struct {
	mu          sync.Mutex
	cleanCtx    context.Context
	cleanCancel context.CancelFunc

	capacity int
	items    *list.List
	itemIdx  map[K]*list.Element
	ttl      time.Duration
}

// NewCache returns a lru cache with given cache size and cache item ttl
func NewCache[K constraints.Ordered, V any](cacheSize int, cacheItemTtl time.Duration) (*Cache[K, V], error) {
	if cacheSize <= 0 {
		return nil, fmt.Errorf("invalid cache size, must be greater than 0")
	}
	if cacheItemTtl <= 0 {
		return nil, fmt.Errorf("invalid cache item ttl, must be greater than 0")
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache[K, V]{
		capacity: cacheSize,
		items:    list.New(),
		itemIdx:  make(map[K]*list.Element),
		ttl:      cacheItemTtl,

		cleanCtx:    ctx,
		cleanCancel: cancel,
	}
	go clean(c)
	return c, nil
}

const cleanInterval = 10 * time.Second

func clean[K constraints.Ordered, V any](c *Cache[K, V]) {
	ontick := func(tick time.Time) {
		ctx, cancel := context.WithTimeout(c.cleanCtx, cleanInterval)
		defer cancel()

		c.mu.Lock()
		defer c.mu.Unlock()

		for e := c.items.Back(); e != nil && ctx.Err() == nil; e = e.Prev() {
			item := e.Value.(*node[K, V])
			if time.Since(item.usedAt) >= c.ttl {
				c.items.Remove(e)
				delete(c.itemIdx, item.key)
			}
		}
	}

	ticker := time.NewTicker(cleanInterval)
	for {
		select {
		case tick := <-ticker.C:
			ontick(tick)
		case <-c.cleanCtx.Done():
			return
		}
	}
}

func exists[K constraints.Ordered, V any](key K, c *Cache[K, V]) (*list.Element, bool) {
	i, ok := c.itemIdx[key]
	return i, ok
}

// PauseCleaning pauses the cleaning of cache items based on the ttl
func (c *Cache[K, T]) PauseCleaning() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cleanCancel()
}

// ResumeCleaning resumes the cleaning of cache items based on the ttl
func (c *Cache[K, T]) ResumeCleaning() {
	if c.cleanCtx.Err() == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	c.cleanCtx = ctx
	c.cleanCancel = cancel

	go clean(c)
}

// Get returns the value and existence of a given key k
func (c *Cache[K, T]) Get(key K) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, doesexist := exists(key, c); doesexist {
		item := el.Value.(*node[K, T])
		c.items.Remove(el)
		delete(c.itemIdx, item.key)
		if time.Since(item.usedAt) >= c.ttl && c.cleanCtx.Err() == nil {
			var zero T
			return zero, false
		}
		// update the item's used at to now and add it to the front of items
		item.usedAt = time.Now()
		c.itemIdx[item.key] = c.items.PushFront(item)

		return item.val, true
	}
	var zero T
	return zero, false
}

// Put puts the given k, v in cache
func (c *Cache[K, T]) Put(key K, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, doesexist := exists(key, c); doesexist {
		item := el.Value.(*node[K, T])
		c.items.Remove(el)

		item.usedAt = time.Now()
		item.val = val
		c.itemIdx[item.key] = c.items.PushFront(item)
		return
	}

	// if the capacity is reached the specified limit, remove the last item and add the new item to front
	if c.items.Len() == c.capacity {
		bval := c.items.Remove(c.items.Back())
		delete(c.itemIdx, bval.(*node[K, T]).key)
	}
	item := &node[K, T]{
		key,
		val,
		time.Now(),
	}
	c.itemIdx[item.key] = c.items.PushFront(item)
}

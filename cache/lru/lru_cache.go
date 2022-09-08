package lru

import (
	"container/list"
	"context"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

type node[K constraints.Ordered, V any] struct {
	key    K
	val    V
	usedAt time.Time
}

type Cache[K constraints.Ordered, V any] struct {
	mu          sync.Mutex
	cleanCtx    context.Context
	cleanCancel context.CancelFunc

	capacity int
	items    *list.List
	itemIdx  map[K]*list.Element
	ttl      time.Duration
}

func NewCache[K constraints.Ordered, V any](cacheSize int, cacheItemTtl time.Duration) *Cache[K, V] {
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
	return c
}

const cleanInterval = 10 * time.Second

func clean[K constraints.Ordered, V any](c *Cache[K, V]) {
	ontick := func(tick time.Time) {
		c.mu.Lock()
		defer c.mu.Unlock()
		ctx, cancel := context.WithTimeout(c.cleanCtx, cleanInterval)
		defer cancel()
		for e := c.items.Front(); e != nil && ctx.Err() == nil; e = e.Next() {
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

func (c *Cache[K, T]) PauseCleaning() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cleanCancel()
}

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
	// create new item to put it in the items
	item := &node[K, T]{
		key,
		val,
		time.Now(),
	}
	c.itemIdx[item.key] = c.items.PushFront(item)
}

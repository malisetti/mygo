package lru

import (
	"context"
	"math"
	"runtime"
	"sync"
	"time"
	"utils"

	"golang.org/x/exp/constraints"
)

type item[K constraints.Ordered, V any] struct {
	Key    K
	Val    V
	UsedAt time.Time
}

type LruCache[K constraints.Ordered, V any] struct {
	mu            sync.Mutex
	cleanCtx      context.Context
	cleanCancel   context.CancelFunc
	cleanInterval time.Duration

	cap   int
	items []*item[K, V]
}

func NewCache[K constraints.Ordered, V any](n int, cleanInterval time.Duration) *LruCache[K, V] {
	ctx, cancel := context.WithCancel(context.Background())
	c := &LruCache[K, V]{
		cap:           n,
		items:         make([]*item[K, V], 0, n),
		cleanCtx:      ctx,
		cleanCancel:   cancel,
		cleanInterval: cleanInterval,
	}
	go clean(c)
	return c
}

func clean[K constraints.Ordered, V any](c *LruCache[K, V]) {
	ontick := func(tick time.Time) {
		c.mu.Lock()
		defer c.mu.Unlock()
		n := len(c.items)
		if n == 0 {
			return
		}
		size := int(math.Ceil(float64(n) / float64(runtime.NumCPU())))
		noSlices := int(math.Ceil(float64(n) / float64(size)))
		newCacheItems := make([][]*item[K, V], noSlices)
		var j int
		var wg sync.WaitGroup
		for i := 0; i < n; i += size {
			j += size
			if j > n {
				j = n
			}
			wg.Add(1)
			go func(x, y int) {
				if c.cleanCtx.Err() == nil {
					return
				}
				defer wg.Done()
				itemSlice := make([]*item[K, V], 0, y-x)
				copy(itemSlice, c.items[x:y])
				for i, n := 0, len(itemSlice); i < n && c.cleanCtx.Err() != nil; i++ {
					e := itemSlice[i]
					if time.Since(e.UsedAt) >= c.cleanInterval {
						itemSlice = removeAt(itemSlice, i)
					}
				}
				if c.cleanCtx.Err() == nil {
					return
				}
				newCacheItems[int(math.Ceil(float64(x)/float64(size)))] = itemSlice
			}(i, j)
		}
		wg.Wait()
		var output []*item[K, V]
		for _, v := range newCacheItems {
			output = append(output, v...) // flatten the 2d items to 1d
		}
		if c.cleanCtx.Err() == nil {
			return
		}
		c.items = output
	}

	ticker := time.NewTicker(c.cleanInterval)
	for {
		select {
		case tick := <-ticker.C:
			ontick(tick)
		case <-c.cleanCtx.Done():
			return
		}
	}
}

func exists[K constraints.Ordered, V any](key K, c *LruCache[K, V]) (int, bool) {
	var compareFunc utils.CompareFunc[*item[K, V]] = func(x *item[K, V]) bool {
		return x.Key == key
	}
	return utils.ExistsAt(c.items, compareFunc)
}

func removeAt[V any](xs []V, i int) []V {
	return append(xs[:i], xs[i+1:]...)
}

func (c *LruCache[K, T]) PauseCleaning() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cleanCancel()
}

func (c *LruCache[K, T]) ResumeCleaning() {
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

func (c *LruCache[K, T]) Get(key K) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.items)
	if n == 0 {
		var zero T
		return zero, false
	}

	if at, doesexist := exists(key, c); doesexist {
		item := c.items[at]
		c.items = removeAt(c.items, at)
		if c.cleanCtx.Err() != nil && time.Since(item.UsedAt) >= c.cleanInterval {
			// remove the item because it is older and cleaning is going on
			var zero T
			return zero, false
		}
		item.UsedAt = time.Now() // update the item's used at to now and add it back to the items
		c.items = append(c.items, item)
		return item.Val, true
	}
	var zero T
	return zero, false
}

func (c *LruCache[K, T]) Put(key K, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if at, doesexist := exists(key, c); doesexist {
		item := c.items[at]
		c.items = removeAt(c.items, at) // remove from index at
		item.UsedAt = time.Now()        // update the used at for the accessed item
		c.items = append(c.items, item) // append the removed item to the items so it becomes the last one
		return
	}

	// create new item to put it in the items
	item := &item[K, T]{
		key,
		val,
		time.Now(),
	}
	if len(c.items) == c.cap { // if the capacity is reached the specified limit, remove the first(zero'th) item and append the new item
		c.items = append(removeAt(c.items, 0), item)
	} else {
		c.items = append(c.items, item) // if capacity is not full, append the item
	}
}

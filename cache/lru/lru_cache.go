package lru

import (
	"context"
	"encoding/json"
	"log"
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
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	cap           int
	items         []*item[K, V]
	cleanInterval time.Duration
}

func NewCache[K constraints.Ordered, V any](n int, cleanInterval time.Duration) *LruCache[K, V] {
	ctx, cancel := context.WithCancel(context.Background())
	c := &LruCache[K, V]{
		cap:           n,
		items:         make([]*item[K, V], 0),
		ctx:           ctx,
		cancel:        cancel,
		cleanInterval: cleanInterval,
	}
	go clean(c)
	return c
}

func clean[K constraints.Ordered, V any](c *LruCache[K, V]) {
	ontick := func(tick time.Time) {
		c.mu.Lock()
		defer c.mu.Unlock()
		if time.Since(tick) >= 999*time.Millisecond {
			log.Println("could not aquire the lock in around a second")
			return
		}

		n := len(c.items)
		if n == 0 {
			return
		}
		size := int(math.Ceil(float64(n) / float64(runtime.NumCPU())))
		cparts := int(math.Ceil(float64(n) / float64(size)))
		newCacheItems := make([][]*item[K, V], cparts)
		var j int
		var wg sync.WaitGroup
		for i := 0; i < n; i += size {
			j += size
			if j > n {
				j = n
			}
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				var cpart []*item[K, V]
				copy(cpart, c.items[x:y])
				for i := 0; i < len(cpart); i++ {
					e := cpart[i]
					at := i
					if time.Since(e.UsedAt) >= c.cleanInterval {
						cpart = append(cpart[:at], cpart[at+1:]...)
					}
				}
				newCacheItems[int(math.Ceil(float64(x)/float64(size)))] = cpart
			}(i, j)
		}
		wg.Wait()
		var output []*item[K, V]
		for _, v := range newCacheItems {
			output = append(output, v...)
		}
		c.items = output
	}

	ticker := time.NewTicker(c.cleanInterval)
	for {
		select {
		case tick := <-ticker.C:
			go ontick(tick)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *LruCache[K, T]) PauseCleaning() {
	c.cancel()
}

func (c *LruCache[K, T]) ResumeCleaning() {
	if c.ctx.Err() == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	go clean(c)
}

func (c *LruCache[K, T]) String() string {
	buf, _ := json.Marshal(c.items)
	return string(buf)
}

func (c *LruCache[K, T]) Get(key K) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.items)
	if n == 0 {
		var zero T
		return zero, false
	}

	at, e := exists(key, c)
	if !e {
		var zero T
		return zero, false
	}
	item := c.items[at]
	c.items = append(c.items[:at], c.items[at+1:]...)
	item.UsedAt = time.Now()
	c.items = append(c.items, item)

	return item.Val, true
}

func exists[K constraints.Ordered, V any](key K, c *LruCache[K, V]) (int, bool) {
	var compareFunc utils.CompareFunc[*item[K, V]] = func(x *item[K, V]) bool {
		return x.Key == key
	}
	return utils.ExistsAt(c.items, compareFunc)
}

func (c *LruCache[K, T]) Put(key K, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	at, e := exists(key, c)
	if e {
		item := c.items[at]
		item.UsedAt = time.Now()
		c.items = append(c.items[:at], c.items[at+1:]...)
		c.items = append(c.items, item)
	} else {
		item := &item[K, T]{
			key,
			val,
			time.Now(),
		}
		if len(c.items) == c.cap {
			c.items = append(c.items[1:], item)
		} else {
			c.items = append(c.items, item)
		}
	}
}

package lru

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"runtime"
	"sync"
	"time"
)

type item[T any] struct {
	Key    string
	Val    T
	UsedAt time.Time
}

type LruCache[T any] struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	cap           int
	items         []*item[T]
	cleanInterval time.Duration
}

func NewCache[T any](n int, interval time.Duration) *LruCache[T] {
	ctx, cancel := context.WithCancel(context.Background())
	c := &LruCache[T]{
		cap:           n,
		items:         make([]*item[T], 0),
		ctx:           ctx,
		cancel:        cancel,
		cleanInterval: interval,
	}
	go clean(c)
	return c
}

const cleanDuration = 1 * time.Second

func clean[T any](c *LruCache[T]) {
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
		newCacheItems := make([][]*item[T], cparts)
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
				var cpart []*item[T]
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
		var output []*item[T]
		for _, v := range newCacheItems {
			output = append(output, v...)
		}
		c.items = output
	}

	ticker := time.NewTicker(cleanDuration)
	for {
		select {
		case tick := <-ticker.C:
			go ontick(tick)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *LruCache[T]) PauseCleaning() {
	c.cancel()
}

func (c *LruCache[T]) ResumeCleaning() {
	if c.ctx.Err() == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	go clean(c)
}

func (c *LruCache[T]) String() string {
	buf, _ := json.Marshal(c.items)
	return string(buf)
}

func (c *LruCache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.items)
	if n == 0 {
		var zero T
		return zero, false
	}

	at, e := exits(key, c)
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

func exits[T any](key string, c *LruCache[T]) (int, bool) {
	n := len(c.items)
	if n == 0 {
		return 0, false
	}
	size := int(math.Ceil(float64(n) / float64(runtime.NumCPU())))
	slices := int(math.Ceil(float64(n) / float64(size)))
	result := make(chan *int, slices)
	var j int
	for i := 0; i < n; i += size {
		j += size
		if j > n {
			j = n
		}
		go func(x, y int) {
			cpart := c.items[x:y]
			for i := 0; i < len(cpart); i++ {
				at := i
				e := cpart[at]
				if e.Key == key {
					exitsAt := x + at
					result <- &exitsAt
					return
				}
			}
			result <- nil
		}(i, j)
	}
	for i := 0; i < slices; i++ {
		v := <-result
		if v != nil {
			return *v, true
		}
	}
	return 0, false
}

func (c *LruCache[T]) Put(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	at, e := exits(key, c)
	if e {
		item := c.items[at]
		item.UsedAt = time.Now()
		c.items = append(c.items[:at], c.items[at+1:]...)
		c.items = append(c.items, item)
	} else {
		item := &item[T]{
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

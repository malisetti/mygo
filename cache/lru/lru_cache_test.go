package lru

import (
	"math/rand"
	"testing"
	"time"
)

type testcase struct {
	key, val string
}

func TestCache(t *testing.T) {
	testcases := []testcase{
		{
			key: "foo",
			val: "bar",
		},
		{
			key: "john",
			val: "doe",
		},
		{
			key: "john1",
			val: "doe1",
		},
		{
			key: "john2",
			val: "doe2",
		},
		{
			key: "john3",
			val: "doe3",
		},
	}
	lruCache := NewCache[string, string](5, 2*time.Second)
	t.Run("test cache gets", func(t *testing.T) {
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			lruCache.Put(testcase.key, testcase.val)
		}
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			val, ok := lruCache.Get(testcase.key)
			if !ok {
				t.Errorf("could not get key %s", testcase.key)
			}
			if val != testcase.val {
				t.Errorf("wanted %s but got %s", testcase.val, val)
			}
		}
	})

	t.Run("test cache size", func(t *testing.T) {
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			lruCache.Put(testcase.key, testcase.val)
		}
		lruCache.Put("john4", "doe4")
		_, ok := lruCache.Get("foo")
		if ok {
			t.Errorf("key \"%s\" should not be present", "foo")
		}
		_, ok = lruCache.Get("john")
		if !ok {
			t.Errorf("key \"%s\" should be present", "john")
		}
	})

	t.Run("test cache ttl", func(t *testing.T) {
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			lruCache.Put(testcase.key, testcase.val)
		}
		time.Sleep(3 * time.Second) // cache item ttl is 2secs
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			_, ok := lruCache.Get(testcase.key)
			if ok {
				t.Errorf("key \"%s\" should not be present", testcase.key)
			}
		}
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BenchmarkCache(b *testing.B) {
	b.Run("bench cache gets", func(b *testing.B) {
		testcases := []testcase{
			{
				key: "foo",
				val: "bar",
			},
			{
				key: "foo1",
				val: "bar1",
			},
			{
				key: "foo2",
				val: "bar2",
			},
			{
				key: "foo3",
				val: "bar3",
			},
		}
		lruCache := NewCache[string, string](100000, 20*time.Second)
		start := time.Now()
		lruCache.Put("fizz", "buzz")
		for _, tc := range testcases {
			lruCache.Put(tc.key, tc.val)
		}
		rand.Seed(time.Now().UnixMilli())
		for i := 0; i < 10000; i++ {
			k, v := randSeq(5), randSeq(5)
			lruCache.Put(k, v)
			if i%5000 == 0 {
				lruCache.Put("foo", "bar")
			}
		}
		var j int
		for i := 0; i < 100000; i++ {
			_, _ = lruCache.Get(testcases[j].key)
			j++
			j = j % len(testcases)
		}

		elapsed := time.Since(start)
		if elapsed < 20*time.Second {

			j = 0
			for i := 0; i < b.N; i++ {
				val, ok := lruCache.Get(testcases[j].key)
				if !ok {
					b.Errorf("key \"%s\" should be present", testcases[j].key)
				}
				if val != testcases[j].val {
					b.Errorf("wanted %s but got \"%s\" for key %s", testcases[j].val, val, testcases[j].key)
				}
				j++
				j = j % len(testcases)
			}
		}
		time.Sleep(20 * time.Second)
		if _, ok := lruCache.Get("fizz"); ok {
			b.Errorf("key \"%s\" should not be present", "fizz")
		}
		lruCache.Put("foo1", "bar1")
		val, ok := lruCache.Get("foo1")
		if !ok {
			b.Errorf("key \"%s\" should be present", "foo1")
		}
		if val != "bar1" {
			b.Errorf("wanted %s but got %s for key %s", "bar1", val, "foo1")
		}
		_, ok = lruCache.Get("foo")
		if ok {
			b.Errorf("key \"%s\" should not be present", "foo")
		}
	})
}

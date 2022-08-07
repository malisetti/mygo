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
	lruCache := NewCache[string, string](1000, 3*time.Second)
	lruCache.Put("foo", "bar")
	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().UnixMilli())
		k, v := randSeq(5), randSeq(5)
		lruCache.Put(k, v)
	}
	lruCache.Put("foo1", "bar1")
	val, ok := lruCache.Get("foo1")
	if !ok {
		b.Errorf("key \"%s\" should be present", "foo1")
	}
	if val != "bar1" {
		b.Errorf("wanted %s but got %s for key %s", "bar1", val, "foo1")
	}
	time.Sleep(3 * time.Second)
	_, ok = lruCache.Get("foo")
	if ok {
		b.Errorf("key \"%s\" should not be present", "foo")
	}
	_, ok = lruCache.Get("foo1")
	if ok {
		b.Errorf("key \"%s\" should not be present", "foo1")
	}
}

package lru

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	type testcase struct {
		key, val string
	}
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
	cache := NewCache[string, string](5, 2*time.Second, 1*time.Second)
	t.Run("test cache gets", func(t *testing.T) {
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			cache.Put(testcase.key, testcase.val)
		}
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			val, ok := cache.Get(testcase.key)
			if !ok {
				t.Errorf("could not get key %s", testcase.key)
			}
			if val != testcase.val {
				t.Errorf("wanted %s but got %s", testcase.val, val)
			}
		}
		time.Sleep(3 * time.Second)
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			_, ok := cache.Get(testcase.key)
			if ok {
				t.Errorf("key \"%s\" should not be present", testcase.key)
			}
		}
	})

	t.Run("test cache size", func(t *testing.T) {
		for i := 0; i < len(testcases); i++ {
			testcase := testcases[i]
			cache.Put(testcase.key, testcase.val)
		}
		cache.Put("john4", "doe4")
		_, ok := cache.Get("foo")
		if ok {
			t.Errorf("key \"%s\" should not be present", "foo")
		}
		_, ok = cache.Get("john")
		if !ok {
			t.Errorf("key \"%s\" should be present", "john")
		}
	})
}

func BenchmarkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

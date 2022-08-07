package main

import (
	"cache/lru"
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	cache := lru.NewCache[string, string](10, 10*time.Second)
	cache.Put("foo", "bar")
	cache.Put("john", "doe")
	cache.Put("john1", "doe")
	cache.Put("john2", "doe")
	cache.Put("john3", "doe")
	cache.Put("john4", "doe")
	cache.Put("john5", "doe")

	fmt.Println(cache.Get("foo"))
	fmt.Println(cache.Get("foo"))
	// go func() {
	// 	rand.Seed(time.Now().UnixNano())
	// 	for i := 0; i < 9990; i++ {
	// 		k := randSeq(5)
	// 		v := randSeq(5)
	// 		cache.Put(k, v)
	// 	}
	// }()

	cache.Put("test", "123")
	fmt.Println(cache.Get("john")) // -1

	fmt.Println(cache.Get("foo"))
	cache.Put("city", "Blore")
	fmt.Println(cache.Get("foo")) // -1

	time.Sleep(12 * time.Second)
	fmt.Println(cache)
	fmt.Println(cache.Get("test")) // -1
	fmt.Println(cache.Get("city")) // -1
}

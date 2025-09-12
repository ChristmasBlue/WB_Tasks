package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu    sync.Mutex
	cache map[int]int
}

func (c *Cache) Add(key int, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[int]int),
	}
}

func main() {
	var wg sync.WaitGroup
	cache := NewCache()
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Worker(cache, &wg, i*10)
	}
	wg.Wait()
	var count int
	for k, v := range cache.cache {
		fmt.Printf("key %d value %d\n", k, v)
		count++
	}
	fmt.Println(count)
}

func Worker(c *Cache, wg *sync.WaitGroup, key int) {
	fmt.Printf("Воркер %d начал работу\n", key/10)
	defer wg.Done()
	defer fmt.Printf("Воркер %d закончил работу\n", key/10)
	for i := key; i < key+5; i++ {
		c.Add(i, i)
	}
}

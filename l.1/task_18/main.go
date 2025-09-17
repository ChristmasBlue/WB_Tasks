package main

import (
	"fmt"
	"sync"
)

type Count struct {
	mu    sync.Mutex
	count int
}

func (c *Count) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func createCount() *Count {
	return &Count{
		count: 0,
	}
}

func main() {
	count := createCount()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count.Increment()
		}()
	}
	wg.Wait()
	fmt.Println(count.count)
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = mergeChannels

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}

// вариант с рекурсией
func mergeChannels(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	switch len(channels) {
	case 0:
		close(out)
	case 1:
		go func() {
			<-channels[0]
			close(out)
		}()

	default:
		go func() {
			defer close(out)
			middle := len(channels) / 2
			select {
			case <-mergeChannels(channels[:middle]...):
			case <-mergeChannels(channels[middle:]...):
			}
		}()
	}

	return out
}

// вариант без рекурсии, с импользование sync.Once{}
func mergeCHannels(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	if len(channels) == 0 {
		close(out)
		return out
	}

	once := sync.Once{}

	for _, channel := range channels {
		go func(ch <-chan interface{}) {
			<-ch
			once.Do(func() {
				close(out)
			})
		}(channel)
	}

	return out
}

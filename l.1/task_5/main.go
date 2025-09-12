package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var t int
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int)
	fmt.Print("Введите время в секундах :")
	fmt.Scan(&t)
	wg.Add(2)
	go WorkerAdd(ctx, &wg, ch)
	go WorkerOut(ctx, &wg, ch)
	<-time.After(time.Duration(t) * time.Second)
	cancel()
	close(ch)
	wg.Wait()
}

func WorkerAdd(ctx context.Context, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	fmt.Println("Воркер загрузчик начал работу")
	i := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Воркер загрузчик закончил работу")
			return
		default:
			ch <- i
			time.Sleep(1 * time.Second)
		}
		i++
	}
}

func WorkerOut(ctx context.Context, wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()
	fmt.Println("Воркер обработчик начал работу")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Воркер обработчик закончил работу")
			return
		default:
			item, ok := <-ch
			if ok {
				fmt.Println(item)
			}
		}
	}
}

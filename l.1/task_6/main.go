package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	ch := make(chan struct{})
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT)
	wg.Add(6)
	go Worker1(ch, &wg)
	go Worker2(ch, &wg)
	go Worker3(ctx, &wg)
	go Worker4(&wg)
	go Worker5(&wg)
	go Worker6(&wg)
	<-signCh
	cancel()
	wg.Wait()
}

// завершение горутины по условию
func Worker1(ch chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	fmt.Println("Воркер 1 начал работу")
	for {
		if i > 10 {
			fmt.Println("Воркер 1 закончил работу")
			close(ch)
			return
		}
		fmt.Println(i)
		i++
		time.Sleep(1 * time.Second)
	}
}

// завершение горутины через канал уведомления
func Worker2(ch <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Воркер 2 начал работу")
	for {
		select {
		case <-ch:
			fmt.Println("Воркер 2 закончил работу")
			return
		default:
			fmt.Println("@")
			time.Sleep(1 * time.Second)
		}
	}
}

// завершение горутины через контекст
func Worker3(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Воркер 3 начал работу")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Воркер 3 закончил работу")
			return
		default:
			fmt.Println("!")
			time.Sleep(1 * time.Second)
		}
	}
}

// завершение горутины естественным образом
func Worker4(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Воркер 4 начал работу")
	time.Sleep(5 * time.Second) //имитация работы
	fmt.Println("Воркер 4 закончил работу")
}

// завершение горутины через runtime.Goexit
func Worker5(wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("Воркер 5 закончил работу")
	fmt.Println("Воркер 5 начал работу")
	i := 10
	for {
		if i == 15 {
			runtime.Goexit()
		}
		fmt.Println(i)
		i++
	}
}

// завершение горутины через восстановление после паники
func Worker6(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			wg.Done()
			fmt.Println("Воркер 6 закончил работу. ", r)
		}
	}()
	fmt.Println("Воркер 6 начал работу")
	time.Sleep(6 * time.Second)
	panic("Конец работы через панику")
}

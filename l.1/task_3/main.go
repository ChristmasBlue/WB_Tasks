package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var workers int
	var wg sync.WaitGroup
	var str string
	ch := make(chan rune, 10)
	fmt.Print("Введите количество воркеров :")
	fmt.Scan(&workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go Worker(&wg, ch)
	}
	fmt.Print("Введите сообщение: ")
	fmt.Scan(&str)
	for _, s := range str {
		ch <- s
	}
	close(ch)
	wg.Wait()
}

func Worker(wg *sync.WaitGroup, ch <-chan rune) {
	defer wg.Done()
	for s := range ch {
		fmt.Printf("%c\n", s)
		time.Sleep(5 * time.Second) //имитация работы горутины
	}
}

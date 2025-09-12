package main

import (
	"fmt"
	"sync"
)

func main() {
	arr := make([]int, 0, 100)
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		arr = append(arr, i)
	}
	chAdd := make(chan int)
	chGet := make(chan int)
	wg.Add(2)
	go Worker1(chAdd, chGet, &wg)
	go Worker2(chGet, &wg)
	for _, v := range arr {
		chAdd <- v
	}
	close(chAdd)
	wg.Wait()
}

func Worker1(chAdd <-chan int, chGet chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Воркер 1 начал работу")
	for item, ok := <-chAdd; ok; item, ok = <-chAdd {
		chGet <- item * 2
	}
	fmt.Println("Воркер 1 закончил работу")
	close(chGet)
}

func Worker2(chGet <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Воркер 2 начал работу")
	for item, ok := <-chGet; ok; item, ok = <-chGet {
		fmt.Println(item)
	}
	fmt.Println("Воркер 2 закончил работу")
}

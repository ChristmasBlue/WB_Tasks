package main

import (
	"fmt"
	"reflect"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan interface{})
	str := "hi"               //string
	num := 5                  //int
	ok := true                //bool
	num1 := 1.5               //float64
	chFake1 := make(chan int) //chan

	wg.Add(1)
	go Worker(ch, &wg)
	ch <- chFake1
	ch <- str
	ch <- num
	ch <- ok
	ch <- num1
	close(ch)
	wg.Wait()
	close(chFake1)
}

func Worker(ch <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for v, ok := <-ch; ok; v, ok = <-ch {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			fmt.Println("Переменная типа: int")
		case reflect.String:
			fmt.Println("Переменная типа: string")
		case reflect.Bool:
			fmt.Println("Переменная типа: bool")
		case reflect.Chan:
			fmt.Println("Переменная типа: chan")
		default:
			fmt.Println("Переменная неизвестного типа")
		}
	}
}

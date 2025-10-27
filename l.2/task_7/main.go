/*package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	//rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}*/

//выведет в случайном порядке числа 1,2,3,4,5,6,7,8 без пробелов, в одну строку
//
//в функции asChan создаётся синхронизированный канал в который пишутся переданные значения,
//в функции merge создаётся синхронизированный канал в который пишутся данные,
//происходит чтение из каналов с помощью select и в case пишутся прочитанные данные из каналов в другой канал,
//как только канал закрывается, в case этому каналу присваивается значение nil, чтобы запретить чтение из канала, и заблокируется case,
//как только оба канала равны nil, третий канал в который происходила запись тоже закрывается,
//в main чтение из канала происходит с помощью range, до тех пор пока канал открыт
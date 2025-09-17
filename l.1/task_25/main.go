package main

import (
	"fmt"
	"time"
)

func main() {
	var sec int
	fmt.Print("Введите длительность(в секундах) блокировки программы: ")
	fmt.Scan(&sec)
	//сначала использую реализацию через цикл
	sleep(time.Duration(sec))
	//второй раз используется реализация через канал+горутина
	ch := make(chan struct{})
	go func() {
		hour, min, seconds := time.Now().Clock()
		fmt.Printf("Время начала блокировки(через канал+горутина): %02d:%02d:%02d\n", hour, min, seconds)
		<-time.After(time.Duration(sec) * time.Second)
		ch <- struct{}{}
	}()
	<-ch
	hour, min, seconds := time.Now().Clock()
	fmt.Printf("Время окончания блокировки(через канал+горутина): %02d:%02d:%02d\n", hour, min, seconds)
	close(ch)
}

// sleep функция которая сравнивает часы, минуты, секунды для выхода из цикла,
// высчитывается конечное время, как только время сравняется цикл прервётся
func sleep(sec time.Duration) {
	now := time.Now()
	hour, min, seconds := now.Clock()
	fmt.Printf("Время начала блокировки(через цикл): %02d:%02d:%02d\n", hour, min, seconds)
	stopHour, stopMin, stopSec := now.Add(sec * time.Second).Clock()
	for {
		if hour, min, seconds = time.Now().Clock(); stopHour == hour && stopMin == min && stopSec == seconds {
			fmt.Printf("Время окончания блокировки(через цикл): %02d:%02d:%02d\n", hour, min, seconds)
			break
		}
	}
}

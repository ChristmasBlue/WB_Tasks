package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	//создаю пустой контекст с функцией отмены,
	//при вызове функции отмены контекст даёт сигнал всем функциям и горутинам в которые он передавался, о завершении,
	//нужно только сделать обработку при получении такого сигнала.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//я создаю канал который слушает системные команды.
	sigCh := make(chan os.Signal, 1)
	//указываю при какой команде этот канал получит сообщение
	signal.Notify(sigCh, syscall.SIGINT)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		fmt.Printf("Воркер %d начал работу\n", i)
		go Worker(ctx, i, &wg)
	}
	<-sigCh
	cancel()
	wg.Wait()
}

func Worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Воркер %d закончил работу\n", id)
			return
		default:
			fmt.Print("@")
			time.Sleep(1 * time.Second)
		}
	}
}

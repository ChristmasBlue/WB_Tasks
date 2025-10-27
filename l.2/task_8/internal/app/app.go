package app

import (
	"context"
	"log"
	"sync"
	"task/internal/worker"
	"task/pkg/console"
)

// Run запуск приложения
func Run() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	hand := worker.NewWorker(ctx, &wg, "0.beevik-ntp.pool.ntp.org")
	go console.StopProgram(ctx, cancel)
	wg.Add(1)
	go hand.GetTime()
	wg.Wait()
	log.Printf("Program stopped, with code: %d\n", hand.StatusCode)
}

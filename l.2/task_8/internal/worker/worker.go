package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"task/internal/domain"
	"task/pkg/console"
	"time"

	"github.com/beevik/ntp"
)

type Worker struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	StatusCode domain.StoppedProg
	server     string
}

// NewWorker конструктор
func NewWorker(ctx context.Context, wg *sync.WaitGroup, server string) *Worker {
	return &Worker{
		ctx:        ctx,
		wg:         wg,
		server:     server,
		StatusCode: domain.StopWithoutErr,
	}
}

// GetTime функция получения времени
func (w *Worker) GetTime() {
	defer w.wg.Done()

	//используем тикер для запросов на сервер
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker is stopping...")
			return
		case <-ticker.C:
			response, err := ntp.Time(w.server)
			if err != nil {
				log.Printf("Error receiving time: %v\n", err)
				w.StatusCode = domain.StopErrGetTime
				return
			}

			//перед каждым выводом времени очищаем консоль
			err = console.Clear()
			if err != nil {
				log.Printf("Error clening console: %v\n", err)
				w.StatusCode = domain.StopErrClearConsole
				return
			}

			fmt.Print(response)
		}
	}
}

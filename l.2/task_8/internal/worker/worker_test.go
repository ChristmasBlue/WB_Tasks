package worker

import (
	"context"
	"errors"
	"sync"
	"task/internal/domain"
	"testing"
	"time"
)

// TestWorker для тестирования - наследует от Worcker
type TestWorker struct {
	*Worker
	MockTime       time.Time //Мок времени для NTP
	MockTimeErr    error     //Мок ошибки для NTP
	MockClearErr   error     //Мок ошибки для очистки консоли
	IterationCount int       // Счетчик итераций (для проверки тикера)
}

// NewTestWorker создаёт тестовый воркер
func NewTestWorker(ctx context.Context, wg *sync.WaitGroup, server string) *TestWorker {
	realWorker := NewWorker(ctx, wg, server)
	return &TestWorker{
		Worker: realWorker,
	}
}

// Переопределяем GetTime для тестов
func (tw *TestWorker) TestGetTime() {
	defer tw.wg.Done()

	ticker := time.NewTicker(50 * time.Microsecond)
	defer ticker.Stop()

	for {
		select {
		case <-tw.ctx.Done():
			return
		case <-ticker.C:
			tw.IterationCount++ // Считаем итерации

			//Используем мок вместо реального NTP
			if tw.MockTimeErr != nil {
				tw.StatusCode = domain.StopErrGetTime
				return
			}

			//Используем мок времени
			_ = tw.MockTime

			//Используем мок ошибки очистки консоли
			if tw.MockClearErr != nil {
				tw.StatusCode = domain.StopErrClearConsole
				return
			}

			if tw.IterationCount >= 5 {
				return
			}
		}
	}
}

// Тесты
func TestWorkerSuccessfulTimeRetrieval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	//Создаём тестовый воркер с моком времени
	worker := NewTestWorker(ctx, &wg, "test-server.pool.ntp.org")
	worker.MockTime = time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
	worker.MockTimeErr = nil
	worker.MockClearErr = nil

	go worker.TestGetTime()

	wg.Wait()
	if worker.StatusCode != domain.StopWithoutErr {
		t.Errorf("Expected status %d (StopWithoutErr), got %d",
			domain.StopWithoutErr, worker.StatusCode)
	}

	if worker.IterationCount == 0 {
		t.Error("Worker should have performed at least one iteration")
	}
}

// TestWorker_NTPError - тест ошибки NTP сервера
func TestWorkerNTPError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	worker := NewTestWorker(ctx, &wg, "invalid-server")
	worker.MockTimeErr = errors.New("NTP server unvailable")
	worker.MockClearErr = nil

	wg.Add(1)
	go worker.TestGetTime()

	wg.Wait()

	if worker.StatusCode != domain.StopErrGetTime {
		t.Errorf("Expected status %d (StopErrGetTime), got %d",
			domain.StopErrGetTime, worker.StatusCode)
	}
}

// TestWorker_ContextCancellation - тест отмены через контекст
func TestWorker_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	worker := NewTestWorker(ctx, &wg, "test-server.pool.ntp.org")
	worker.MockTime = time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)

	// Быстро отменяем контекст
	go worker.TestGetTime()
	time.Sleep(20 * time.Millisecond) // Даем немного времени на старт
	cancel()

	wg.Wait()

	// При отмене контекста статус должен остаться StopWithoutErr
	if worker.StatusCode != domain.StopWithoutErr {
		t.Errorf("Expected status %d (StopWithoutErr) on context cancellation, got %d",
			domain.StopWithoutErr, worker.StatusCode)
	}
}

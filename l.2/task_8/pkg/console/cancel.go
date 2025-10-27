package console

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func StopProgram(ctx context.Context, cancel context.CancelFunc) {
	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT)

	<-signCh
	cancel()
}

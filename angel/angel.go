package angel

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Interuptor(cancel context.CancelFunc) {
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)
		select {
		case <-s:
			cancel()
		}
	}()
}

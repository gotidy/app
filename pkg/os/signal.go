package os

import (
	"context"
	"os/signal"
	"syscall"
)

// ContextWithSignal creates a context canceled when SIGINT or SIGTERM are notified.
func ContextWithSignal(ctx context.Context) context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		stop()
	}()

	return ctx
}

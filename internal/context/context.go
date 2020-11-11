package context

import (
	"context"
	"sync"

	"github.com/gotidy/app/internal/log"
	"github.com/gotidy/app/pkg/os"
)

// Context intended to be passed to a command specific Run function.
type Context struct {
	// WaitGroup is used to wait for asynchronous workers to complete.
	WaitGroup *sync.WaitGroup
	// Ctx will be canceled when SIGINT or SIGTERM are notified.
	Ctx context.Context

	// Some parameters
	Logger *log.Logger
}

func NewContext() *Context {
	return &Context{
		Logger:    log.New(),
		Ctx:       os.ContextWithSignal(context.Background()),
		WaitGroup: &sync.WaitGroup{},
	}
}

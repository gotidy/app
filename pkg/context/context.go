package context

import (
	"context"
	"os"
	"sync"

	"github.com/rs/zerolog"

	signal "github.com/gotidy/app/pkg/os"
)

// Fields is alias to zerolog.Context.
type Fields = zerolog.Context

// Context of execution. It consists of logger and cancelation context.
type Context struct {
	// WaitGroup is used to wait for asynchronous workers to complete.
	WaitGroup *sync.WaitGroup
	// Ctx will be canceled when SIGINT or SIGTERM are notified.
	Ctx context.Context

	// Some parameters
	Logger zerolog.Logger
}

// New create new Context.
func New() Context {
	return Context{
		Logger:    zerolog.New(os.Stderr).With().Timestamp().Logger(),
		Ctx:       signal.ContextWithSignal(context.Background()),
		WaitGroup: &sync.WaitGroup{},
	}
}

// WithFields sets Logger and returns a copy of the context.
func (c Context) WithLogger(l zerolog.Logger) Context {
	c.Logger = l
	return c
}

// WithFields sets Logger fields and returns a copy of the context.
//
//   ctx := context.New().WithFields(func(f context.Fields) context.Fields){
//       return f.Str("Version", "v1.2.3")
//   })
func (c Context) WithFields(f func(f Fields) Fields) Context {
	c.Logger = f(c.Logger.With()).Logger()
	return c
}

// WithContext sets Context and returns self copy.
func (c Context) WithContext(ctx context.Context) Context {
	c.Ctx = ctx
	return c
}

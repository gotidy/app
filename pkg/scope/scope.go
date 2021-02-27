package scope

import (
	"context"
	"os"
	"sync"

	"github.com/rs/zerolog"

	signal "github.com/gotidy/app/pkg/os"
)

// Fields is alias to zerolog.Context.
type Fields = zerolog.Context

// Scope of execution. It consists of logger and cancelation context.
type Scope struct {
	// WaitGroup is used to wait for asynchronous workers to complete.
	WaitGroup *sync.WaitGroup
	// Ctx will be canceled when SIGINT or SIGTERM are notified.
	Ctx context.Context

	// Some parameters
	Logger zerolog.Logger
}

// New create new Scope.
func New() Scope {
	return Scope{
		Logger:    zerolog.New(os.Stderr).With().Timestamp().Logger(),
		Ctx:       signal.ContextWithSignal(context.Background()),
		WaitGroup: &sync.WaitGroup{},
	}
}

// WithFields sets Logger and returns a copy of the context.
func (c Scope) WithLogger(l zerolog.Logger) Scope {
	c.Logger = l
	return c
}

// WithFields sets Logger fields and returns a copy of the context.
//
//   ctx := scope.New().WithFields(func(f scope.Fields) scope.Fields){
//       return f.Str("Version", "v1.2.3")
//   })
func (c Scope) WithFields(f func(f Fields) Fields) Scope {
	c.Logger = f(c.Logger.With()).Logger()
	return c
}

func (c Scope) WithRole(name string) Scope {
	c.Logger = c.Logger.With().Str("role", name).Logger()
	return c
}

// WithContext sets Context and returns self copy.
func (c Scope) WithContext(ctx context.Context) Scope {
	c.Ctx = ctx
	return c
}

package server

import (
	"github.com/rs/zerolog"
)

type Handlers struct {
	Healthcheck
}

func NewHandlers(log zerolog.Logger) *Handlers {
	return &Handlers{
		Healthcheck: NewHealthcheck(log),
	}
}

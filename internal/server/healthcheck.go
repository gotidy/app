package server

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Healthcheck struct {
	Base
}

func NewHealthcheck(log zerolog.Logger) Healthcheck {
	return Healthcheck{Base: NewBase(log)}
}

func (h Healthcheck) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	// h.ok(w)
	w.WriteHeader(http.StatusOK)
}

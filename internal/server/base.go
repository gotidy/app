package server

import (
	"github.com/rs/zerolog"
)

type Base struct {
	log zerolog.Logger
}

func NewBase(log zerolog.Logger) Base {
	// return Base{log: scope.Logger.With().Str("context", "http-server").Logger()}
	return Base{log: log}
}

type Raw []byte

// func (b Base) response(w http.ResponseWriter, status int, body ...interface{}) {
// 	w.WriteHeader(status)
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")

// 	var err error
// 	if len(body) > 0 {
// 		if bytes, ok := body[0].(Raw); ok {
// 			_, err = w.Write(bytes)
// 		} else {
// 			err = json.NewEncoder(w).Encode(body)
// 		}
// 	}
// 	if err != nil {
// 		b.log.Err(err).Msg("response")
// 	}
// 	b.log.Debug().Int("status", status).Msg("response")
// }

// func (b Base) ok(w http.ResponseWriter, body ...interface{}) {
// 	b.response(w, http.StatusOK, body...)
// }

// func (b Base) error(w http.ResponseWriter, err error) {
// 	status := http.StatusInternalServerError
// 	var body interface{}
// 	if err, ok := err.(*HTTPError); ok {
// 		body = err.ResponseBody()
// 	}

// 	b.response(w, status, body)
// }

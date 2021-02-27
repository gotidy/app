package hlog

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Middleware func(http.Handler) http.Handler

type ChainHandler struct {
	Middlewares []Middleware
}

func Chain(m ...Middleware) ChainHandler {
	return ChainHandler{Middlewares: m}
}

func (c ChainHandler) Use(m ...Middleware) ChainHandler {
	c.Middlewares = append(c.Middlewares, m...)
	return c
}

func (c ChainHandler) Handle(h http.Handler) http.Handler {
	for i := range c.Middlewares {
		h = c.Middlewares[len(c.Middlewares)-1-i](h)
	}
	return h
}

func DumpHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dump, err := httputil.DumpRequest(r, true); err == nil {
			hlog.FromRequest(r).Debug().
				Str("dump", string(dump)).
				Msg("request dump")
		}
		next.ServeHTTP(w, r)
	})
}

func Handler(logger zerolog.Logger, next http.Handler) http.Handler {
	return Chain().
		Use(middleware.Recoverer).
		Use(hlog.NewHandler(logger)).
		Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Debug().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Str("operation", "request").
				Msg("processing request")
		})).
		Use(hlog.RemoteAddrHandler("ip")).
		// If you are service is behind load balancer like nginx, you might want to
		// use X-Request-ID instead of injecting request id. You can do some thing
		// like this,
		// r.Use(hlog.CustomHeaderHandler("reqId", "X-Request-Id"))
		Use(hlog.RequestIDHandler("req_id", "Request-Id")).
		// Use(DumpHandler).
		Handle(next)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientLogger struct {
	log    zerolog.Logger
	client HTTPClient
}

func Client(log zerolog.Logger, client HTTPClient) *ClientLogger {
	return &ClientLogger{
		log:    log,
		client: client,
	}
}

func (c ClientLogger) Do(r *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := c.client.Do(r)
	latency := time.Since(start)

	var event *zerolog.Event
	if err != nil {
		event = c.log.Error().Err(err)
	} else {
		event = c.log.Debug()
	}
	event = event.Str("method", r.Method).Str("url", r.URL.String()).Dur("latency", latency)

	if res != nil {
		event.Int("status", res.StatusCode).
			Int64("size", res.ContentLength)
	} else {
		event.Int("status", -1).
			Int64("size", -1)
	}

	event.Msg("sending a request")

	return res, err
}

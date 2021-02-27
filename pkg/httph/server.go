// HTTP helpers
package httph

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gotidy/app/pkg/log/hlog"
	"github.com/gotidy/app/pkg/scope"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	Server  *http.Server
	Addr    *url.URL
	Metrics bool
}

type Option func(opts *server)

// WithServer set the custom http.Server, else is used default.
func WithServer(srv *http.Server) Option {
	return func(opts *server) {
		opts.Server = srv
	}
}

// WithMetrics adds metrics at the "/metrics" path.
func WithMetrics() Option {
	return func(opts *server) {
		opts.Metrics = true
	}
}

// WithAddr optionally specifies the TCP address for the server to listen on, in the form "host:port".
// If empty, ":http" (port 80) is used. The service names are defined in RFC 6335 and assigned by IANA.
// See net.Dial for details of the address format.
func WithAddr(addr *url.URL) Option {
	return func(opts *server) {
		opts.Addr = addr
	}
}

// Serve start server with logging requests and metrics.
func Serve(scope scope.Scope, handler http.Handler, opts ...Option) {
	srv := &server{}
	for _, opt := range opts {
		opt(srv)
	}
	if srv.Server == nil {
		srv.Server = &http.Server{}
	}
	if srv.Addr == nil {
		srv.Addr = &url.URL{Host: ":80", Scheme: "http", Path: "/"}
	}
	if srv.Addr.Path == "" {
		srv.Addr.Path = "/"
	}

	mux := http.NewServeMux()
	if srv.Metrics {
		mux.Handle("/metrics", promhttp.Handler())
		handler = metrics(
			handler,
		)
	}
	mux.Handle(srv.Addr.Path, hlog.Handler(scope.Logger, handler))

	srv.Server.Handler = mux
	srv.Server.Addr = srv.Addr.Host

	scope.WaitGroup.Add(1)

	log := scope.Logger.With().Str("host", srv.Addr.Hostname()).Str("port", srv.Addr.Port()).Logger()

	go func() {
		<-scope.Ctx.Done()
		log.Info().Msg("Stopping HTTP server")
		if err := srv.Server.Shutdown(context.Background()); err != nil {
			log.Info().Err(err).Msg("HTTP server Shutdown")
		}
		scope.WaitGroup.Done()
	}()

	go func() {
		if err := srv.Server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("HTTP server listening failed")
		}
		log.Info().Msg("HTTP server stopped")
	}()

	log.Info().Msg("HTTP server started")
}

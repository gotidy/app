package server

import (
	"net/http"
	"net/url"

	"github.com/gotidy/app/internal/server/api"
	"github.com/gotidy/app/pkg/httph"
	"github.com/gotidy/app/pkg/scope"
)

func Serve(scope scope.Scope, addr *url.URL) {
	httph.Serve(scope,
		api.Handler(NewHandlers(scope.Logger)),
		httph.WithServer(&http.Server{}),
		httph.WithAddr(addr),
		httph.WithMetrics(),
	)
}

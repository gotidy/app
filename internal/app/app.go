package app

import (
	"net/url"

	"github.com/gotidy/app/internal/server"
	"github.com/gotidy/app/pkg/config"
	"github.com/gotidy/app/pkg/scope"
)

type Config struct {
	DB  config.Connection
	URL *url.URL
}

func Run(scope scope.Scope, conf *Config) error {
	server.Serve(scope.WithRole("http-server"), conf.URL)

	scope.Logger.Info().Msg("Application successfully started")

	return nil
}

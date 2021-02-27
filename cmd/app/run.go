package main

import (
	"net/url"

	"github.com/gotidy/app/internal/app"
	"github.com/gotidy/app/pkg/config"
	"github.com/gotidy/app/pkg/scope"
)

type DB struct {
	Host     string `help:"Database host" hidden:""`
	Port     int    `help:"Database port" hidden:""`
	User     string `help:"Database user" hidden:""`
	Password string `help:"Database password" hidden:""`
	Name     string `help:"Database name" hidden:""`
}

// Run command (default).
type RunCmd struct {
	// DB connection info
	DB `prefix:"db." embed:""`

	Scheme string `default:"http" help:"Scheme is part of «http(s)://<host>:<port>/<path>» that is used as webhook for Aidbox."`
	Host   string `help:"Host is part of «http(s)://<host>:<port>/<path> that is used as webhook for Aidbox."`
	Port   int    `default:"80" help:"Port is part of «http(s)://<host>:<port>/<path> that is used as webhook for Aidbox."`
	Path   string `default:"/" help:"Path is part of «http(s)://<host>:<port>/<path> that is used as webhook for Aidbox."`
}

func (c *RunCmd) Run(scope *scope.Scope) error {
	// Some code that will be to execute command.
	return app.Run(*scope, &app.Config{
		DB: config.Connection{
			Address: config.NetAddress{Host: c.DB.Host, Port: c.DB.Port},
			User:    config.UserCredential{User: c.DB.User, Password: &c.DB.Password},
		},
		URL: &url.URL{Scheme: c.Scheme, Host: config.NetAddress{Host: c.Host, Port: c.Port}.String(), Path: c.Path},
	})
}

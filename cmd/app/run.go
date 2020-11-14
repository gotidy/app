package main

import (
	"github.com/gotidy/app/internal/app"
	"github.com/gotidy/app/pkg/config"
	"github.com/gotidy/app/pkg/context"
)

type DB struct {
	Host     string `help:"Database host."`
	Port     int    `help:"Database port."`
	User     string `help:"Database user."`
	Password string `help:"Database password."`
	Name     string `help:"Database name."`
}

// Run command (default).
type RunCmd struct {
	// DB connection info
	DB `prefix:"db." embed:""`
}

func (c *RunCmd) Run(ctx *context.Context) error {
	// Some code that will be to execute command.
	return app.Run(ctx, &app.Config{
		DB: config.Connection{
			Address: config.NetAddress{Host: c.DB.Host, Port: c.DB.Port},
			User:    config.UserCredential{User: c.DB.User, Password: &c.DB.Password},
		},
	})
}

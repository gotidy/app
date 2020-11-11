package main

import (
	"github.com/gotidy/app/internal/app"
	"github.com/gotidy/app/internal/config"
	"github.com/gotidy/app/internal/context"
)

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
}

// Run command (default)
type RunCmd struct {
	// DB connection info
	DB `prefix:"db." embed:""`
}

func (c *RunCmd) Run(ctx *context.Context) {
	// Some code that will be to execute command.
	app.Run(ctx, &app.Config{
		DB: config.Connection{
			Address: config.NetAddress{Host: c.DB.Host, Port: c.DB.Port},
			User:    config.UserCredential{User: c.DB.User, Password: &c.DB.Password},
		},
	})
}

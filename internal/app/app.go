package app

import (
	"github.com/gotidy/app/pkg/config"
	"github.com/gotidy/app/pkg/context"
)

type Config struct {
	DB config.Connection
}

// type App struct {
// 	config *Config
// }

// func NewApp(conf *Config) *App {
// 	return &App{config: conf}
// }

func Run(ctx *context.Context, conf *Config) error {
	// defer ctx.WaitGroup.Done()

	ctx.Logger.Info().Msg("Application succeful started")

	// Some code

	return nil
}

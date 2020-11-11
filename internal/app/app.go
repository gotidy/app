package app

import (
	"github.com/gotidy/app/internal/config"
	"github.com/gotidy/app/internal/context"
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

func Run(ctx *context.Context, conf *Config) {
	defer ctx.WaitGroup.Done()

	// Some code
}

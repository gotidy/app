package main

import (
	"github.com/gotidy/app/internal/cli"
	"github.com/gotidy/app/internal/context"
)

func main() {
	cli.Run(
		NewCli(),
		context.NewContext(),
		cli.Name(ApplicationName),
		cli.Version(Version),
		cli.Env("APP"),
	)
}

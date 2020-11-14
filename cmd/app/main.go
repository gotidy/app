package main

import (
	"github.com/gotidy/app/pkg/cli"
	"github.com/gotidy/app/pkg/context"
)

func main() {
	cli.Run(
		NewCli(),
		context.New().WithFields(func(f context.Fields) context.Fields {
			return f.Str("Version", Version)
		}),
		cli.Name(ApplicationName),
		cli.Version(ApplicationName+" "+Version),
		cli.Env("APP"),
	)
}

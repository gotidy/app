package main

import (
	"github.com/gotidy/app/pkg/cli"
	"github.com/gotidy/app/pkg/scope"
)

func main() {
	cli.Run(
		NewCli(),
		scope.New().WithFields(func(f scope.Fields) scope.Fields {
			return f.Str("Version", version+"-"+commit)
		}),
		cli.Name(ApplicationName),
		cli.Version(ApplicationName+" "+version),
		cli.Env("APP"),
	)
}

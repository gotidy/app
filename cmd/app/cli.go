package main

type Cli struct {
	// Default command
	Run RunCmd `cmd:"" name:"" help:"Run application" default:"1"`
}

func NewCli() *Cli {
	return &Cli{}
}

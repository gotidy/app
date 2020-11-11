package cli

import (
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/gotidy/app/internal/context"
	"github.com/gotidy/app/internal/log"
)

type LoggingFormat string

func (l LoggingFormat) AfterApply(ctx *context.Context) error {
	// format := string(ctx.FlagValue(trace.Flag).(LoggingFormat))
	switch strings.ToLower(string(l)) {
	case "text":
		ctx.Logger.Format(log.TextFormat)
	case "json":
		ctx.Logger.Format(log.JSONFormat)
	}
	return nil
}

type LoggingLevel string

func (l LoggingLevel) BeforeApply(ctx *context.Context) error {
	switch l {
	case "debug":
		ctx.Logger.Level(log.DebugLevel)
	case "info":
		ctx.Logger.Level(log.InfoLevel)
	case "warn":
		ctx.Logger.Level(log.WarnLevel)
	case "error":
		ctx.Logger.Level(log.ErrorLevel)
	case "fatal":
		ctx.Logger.Level(log.FatalLevel)
	case "panic":
		ctx.Logger.Level(log.PanicLevel)
	case "trace":
		ctx.Logger.Level(log.TraceLevel)
	}
	return nil
}

type LoggingOutput string

func (l LoggingOutput) AfterApply(ctx *context.Context) (err error) {
	// out := string(ctx.FlagValue(trace.Flag).(LoggingOutput))
	out := string(l)
	var w io.Writer
	switch strings.ToUpper(string(out)) {
	case "", "STDERR":
		w = os.Stderr
	case "STDOUT":
		w = os.Stdout
	default:
		if w, err = os.OpenFile(string(out), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			return err
		}
	}
	ctx.Logger.Output(w)
	return nil
}

// type ConfigFlag string

// // BeforeResolve adds a resolver.
// func (c ConfigFlag) BeforeResolve(kong *kong.Kong, ctx *kong.Context, trace *kong.Path) error {
// 	path := string(ctx.FlagValue(trace.Flag).(ConfigFlag))
// 	if path == "" {
// 		return nil
// 	}
// 	resolver, err := kong.LoadConfig(path)
// 	if err != nil {
// 		return err
// 	}
// 	ctx.AddResolver(resolver)
// 	return nil
// }

type Cli struct {
	// Common flags.
	LoggingFormat LoggingFormat `help:"Logging format (${enum}). Default value: \"${default}\"" enum:"text,json" default:"text"`
	LoggingLevel  LoggingLevel  `help:"Logging level (${enum})." enum:"debug,info,warn,error,fatal,panic,trace" default:"info"`
	LoggingOutput LoggingOutput `help:"Logging output (stderr,stdout,<path>)." default:"stderr"`

	// if the flag --config=<path> is defined, then config will be loaded
	Config kong.ConfigFlag `type:"path" help:"Config path."`
	// Used for showing version if defined --version flag
	Version kong.VersionFlag
}

type Option func() kong.Option

func Name(name string) Option {
	return func() kong.Option { return kong.Name(name) }
}

func Version(version string) Option {
	return func() kong.Option { return kong.Vars{"version": version} }
}

func Paths(paths ...string) Option {
	return func() kong.Option { return kong.Configuration(JSON, paths...) }
}

func Run(cli interface{}, ctx *context.Context, options ...Option) {
	root, err := CombineStructs(Cli{}, cli)
	if err != nil {
		ctx.Logger.ApplicationStartFailed(err)
	}

	var kongOptions []kong.Option
	for _, option := range options {
		kongOptions = append(kongOptions, option())
	}

	kongCtx := kong.Parse(root, kongOptions...)
	if err = CopyStruct(root, cli); err != nil {
		ctx.Logger.ApplicationStartFailed(err)
	}
	err = kongCtx.Run(ctx)
	if err != nil {
		ctx.Logger.ApplicationStartFailed(err)
	}

	ctx.WaitGroup.Wait()
	ctx.Logger.ApplicationSuccefulStopped()
}

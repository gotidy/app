package cli

import (
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/alecthomas/kong"

	"github.com/gotidy/app/pkg/context"
)

type LoggingFormat string

func (l LoggingFormat) AfterApply(ctx *context.Context, writer *ConsoleWriter) error {
	format := ColoredTextFormat
	switch strings.ToLower(string(l)) {
	case "coloredtext":
		format = ColoredTextFormat
	case "text":
		format = TextFormat
	case "json":
		format = JSONFormat
	}
	*ctx = ctx.WithLogger(ctx.Logger.Output(writer.Format(format)))
	return nil
}

type LoggingLevel string

func (l LoggingLevel) BeforeApply(ctx *context.Context) error {
	level := InfoLevel
	switch l {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	case "panic":
		level = PanicLevel
	case "trace":
		level = TraceLevel
	}
	*ctx = ctx.WithLogger(ctx.Logger.Level(level))
	return nil
}

type LoggingOutput string

func (l LoggingOutput) AfterApply(ctx *context.Context, writer *ConsoleWriter) (err error) {
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
	*ctx = ctx.WithLogger(ctx.Logger.Output(writer.Out(w)))
	return nil
}

type Cli struct {
	// Common flags.
	LoggingFormat LoggingFormat `help:"Logging format (${enum}). Default value: \"${default}\"" enum:"json,text,coloredtext" default:"coloredtext"`
	LoggingLevel  LoggingLevel  `help:"Logging level (${enum})." enum:"debug,info,warn,error,fatal,panic,trace" default:"info"`
	LoggingOutput LoggingOutput `help:"Logging output (stderr,stdout,<path>)." default:"stderr"`

	// if the flag --config=<path> is defined, then config will be loaded
	Config kong.ConfigFlag `type:"path" help:"Config path."`
	// Used for showing version if defined --version flag
	Version kong.VersionFlag
}

type Option func() kong.Option

// Name sets application name.
func Name(name string) Option {
	return func() kong.Option { return kong.Name(name) }
}

// Version sets application version. It outputs when --version flag is defined.
func Version(version string) Option {
	return func() kong.Option { return kong.Vars{"version": version} }
}

// Path sets configuration pathes.
func Paths(paths ...string) Option {
	return func() kong.Option { return kong.Configuration(JSON, paths...) }
}

// Env inits environment names for flags.
// For example:
//   --some.value -> PREFIX_SOME_VALUE
func Env(prefix string) Option {
	var f kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		switch env := flag.Env; {
		case flag.Name == "help":
			return nil, nil
		case env == "-":
			flag.Env = ""
			return nil, nil
		case env != "":
			return nil, nil
		}
		replacer := strings.NewReplacer("-", "_", ".", "_")
		name := replacer.Replace(flag.Name)
		// Split by upper chars "SomeOne" -> ["Some", "One"]
		var names []string
		if prefix != "" {
			names = []string{prefix}
		}
		for {
			i := strings.IndexFunc(name, unicode.IsUpper)
			if i < 0 {
				names = append(names, strings.Trim(name, "_"))
				break
			}
			names = append(names, strings.Trim(name[:i], "_"))
			name = name[i:]
		}
		name = strings.ToUpper(strings.Join(names, "_"))
		flag.Env = name
		flag.Value.Tag.Env = name
		return nil, nil
	}

	return func() kong.Option { return kong.Resolvers(f) }
}

// Run executes the Run() method on the selected command, which must exist.
func Run(cli interface{}, ctx context.Context, options ...Option) {
	root, err := CombineStructs(Cli{}, cli)
	if err != nil {
		ApplicationStartFailed(ctx.Logger, err)
	}

	kongOptions := make([]kong.Option, 0, len(options))
	for _, option := range options {
		kongOptions = append(kongOptions, option())
	}

	kongOptions = append(kongOptions, kong.Bind(&ctx), kong.Bind(NewConsoleWriter()), kong.Resolvers())
	kongCtx := kong.Parse(root, kongOptions...)
	if err = CopyStruct(root, cli); err != nil {
		ApplicationStartFailed(ctx.Logger, err)
	}
	err = kongCtx.Run(ctx)
	if err != nil {
		ApplicationStartFailed(ctx.Logger, err)
	}

	ctx.WaitGroup.Wait()
	ApplicationSuccefulStopped(ctx.Logger)
}

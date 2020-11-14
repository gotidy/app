package cli

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Format int

const (
	TextFormat Format = iota
	ColoredTextFormat
	JSONFormat
)

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = zerolog.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = zerolog.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = zerolog.NoLevel
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
	// TraceLevel defines trace log level.
	TraceLevel = zerolog.TraceLevel
)

type ConsoleWriter struct {
	w      zerolog.ConsoleWriter
	format Format
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{
		w: zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
			//FormatLevel: func(i interface{}) string {
			//	return strings.ToUpper(fmt.Sprintf("%-6s|", i))
			//},
			TimeFormat: time.RFC3339,
		},
		format: ColoredTextFormat,
	}
}

func (c *ConsoleWriter) Format(f Format) *ConsoleWriter {
	c.format = f
	c.w.NoColor = f != ColoredTextFormat
	return c
}

func (c *ConsoleWriter) Out(w io.Writer) *ConsoleWriter {
	c.w.Out = w
	return c
}

func (c *ConsoleWriter) Write(p []byte) (n int, err error) {
	if c.format == JSONFormat {
		return c.w.Out.Write(p)
	}
	return c.w.Write(p)
}

package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Level = zerolog.Level

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

type Format int

const (
	TextFormat Format = iota
	ColoredTextFormat
	JSONFormat
)

func applyFormat(f Format, w io.Writer) io.Writer {
	switch f {
	case JSONFormat:
		return w
	default:
		return zerolog.ConsoleWriter{
			Out:             w,
			NoColor:         f != ColoredTextFormat,
			FormatLevel:     func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%-6s|", i)) },
			FormatFieldName: func(i interface{}) string { return fmt.Sprintf("%s=", i) },
		}
	}
}

type Logger struct {
	zerolog.Logger
	out    io.Writer
	format Format
}

func New() *Logger {
	return (&Logger{
		Logger: zerolog.New(os.Stderr),
		format: TextFormat,
	}).Output(os.Stderr)
}

func (l *Logger) Format(f Format) *Logger {
	l.format = f
	return l.Output(l.out)
}

func (l *Logger) Level(lvl Level) *Logger {
	l.Logger = l.Logger.Level(lvl)
	return l
}

func (l *Logger) Output(w io.Writer) *Logger {
	l.Logger = l.Logger.Output(applyFormat(l.format, w))
	return l
}

func (l *Logger) ApplicationSuccefulStopped() {
	l.Logger.Info().Msg("Application succeful stopped")
}

func (l *Logger) ApplicationStartFailed(err error) {
	l.Logger.Fatal().Err(err).Msg("Application start failed")
}

func (l *Logger) ApplicationStarted() {
	l.Logger.Info().Msg("Application succeful started")
}

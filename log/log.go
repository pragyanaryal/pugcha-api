package log

import (
	"context"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

func init() {
	logLevel := zerolog.InfoLevel
	zerolog.SetGlobalLevel(logLevel)
	if config.Configuration.LevelLog == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// os.Stderr
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	Logger = &log
}

func Output(w io.Writer) zerolog.Logger {
	return Logger.Output(w)
}

func With() zerolog.Context {
	return Logger.With()
}

func Level(level zerolog.Level) zerolog.Logger {
	return Logger.Level(level)
}

func Sample(s zerolog.Sampler) zerolog.Logger {
	return Logger.Sample(s)
}

func Err(err error) *zerolog.Event {
	return Logger.Err(err)
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

func Panic() *zerolog.Event {
	return Logger.Panic()
}

func Log() *zerolog.Event {
	return Logger.Log()
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func WithLevel(level zerolog.Level) *zerolog.Event {
	return Logger.WithLevel(level)
}

func Hook(h zerolog.Hook) zerolog.Logger {
	return Logger.Hook(h)
}

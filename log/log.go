package log

import (
	"io"
	"os"
	"reflect"
	"time"

	"github.com/rs/zerolog"
)

const (
	DEBUG_LEVEL = zerolog.DebugLevel
	INFO_LEVEL  = zerolog.InfoLevel
	WARN_LEVEL  = zerolog.WarnLevel
	ERR_LEVEL   = zerolog.ErrorLevel
	FATAL_LEVEL = zerolog.FatalLevel

	KeyModule = "mod"
	KeyEvent  = "event"

	ModuleIndexer = "indexer"
	ModuleCLI     = "cli"
	ModuleServer  = "server"
	ModuleStore   = "store"
	ModuleService = "service"
	ModuleChain   = "chain"
	ModuleJob     = "job"
)

// global
var Logger = zerolog.New(writer).With().Timestamp().Logger().Level(zerolog.ErrorLevel)
var ConsoleWriter = zerolog.ConsoleWriter{Out: os.Stdout}

var (
	DefaultLevel = DEBUG_LEVEL

	writer            = &Writer{Out: os.Stderr}
	defaultWriter     = ConsoleWriter
	defaultTimeFormat = time.RFC3339
)

// default conf
func init() {
	SetLevel(DefaultLevel)
	SetWriter(defaultWriter)
	SetTimeFormat(defaultTimeFormat)
}

func SetLevel(lv zerolog.Level) {
	Logger = Logger.Level(lv)
}

func SetWriter(w io.Writer) {
	writer.SetWriter(w)
}

func SetTimeFormat(format string) {
	zerolog.TimeFieldFormat = format
	if reflect.TypeOf(writer.Out).Name() == "ConsoleWriter" {
		// reset writer
		ConsoleWriter = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: format,
		}
		SetWriter(ConsoleWriter)
	}
}

func SetLogColor(withColor bool) {
	ConsoleWriter.NoColor = !withColor
	SetWriter(ConsoleWriter)
}

func ToStr(lv zerolog.Level) string {
	switch lv {
	case zerolog.DebugLevel:
		return "debug"
	case zerolog.InfoLevel:
		return "info"
	case zerolog.WarnLevel:
		return "warn"
	case zerolog.ErrorLevel:
		return "err"
	case zerolog.FatalLevel:
		return "fatal"
	default:
		panic("unexpected logging level")
	}
}

func ToLogLevel(lv string) zerolog.Level {
	switch lv {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "err":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		panic("unexpected logging level string")
	}
}

func Indexer() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleIndexer).Logger()
}

func CLI() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleCLI).Logger()
}

func Server() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleServer).Logger()
}

func Store() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleStore).Logger()
}

func Service() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleService).Logger()
}

func Chain() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleChain).Logger()
}

func Job() zerolog.Logger {
	return Logger.With().Str(KeyModule, ModuleJob).Logger()
}

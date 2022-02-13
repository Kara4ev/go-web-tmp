package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type level int

const (
	DebugLevel level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

var logger *zerolog.Logger

func init() {
	l := zerolog.New(os.Stdout)
	logger = &l
	setLogLevel("debug")
}

func InitLogger(level, filePath string) {

	err := os.MkdirAll("logs", 0755)

	if err != nil {
		panic("can't create log dir. no configured logging to files")
	}

	fileLog, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic("can't open/create log file. no configured logging to files")
	}

	setLogLevel(level)

	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zerolog.CallerFieldName = "exe"

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, fileLog)
	l := zerolog.New(multi).With().Timestamp().CallerWithSkipFrameCount(5).Logger()

	logger = &l
}

func setLogLevel(level string) {
	var l zerolog.Level

	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	case "fatal":
		l = zerolog.FatalLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)
	logger.Level(l)
}

func msg(level level, message interface{}, args ...interface{}) {

	m := getString(message)
	switch level {
	case InfoLevel:
		logInfo(m, args...)
	case DebugLevel:
		logDebug(m, args...)
	case WarnLevel:
		logWarning(m, args...)
	case ErrorLevel:
		logError(m, args...)
	case FatalLevel:
		logfatal(m, args...)
	}
}

func getString(message interface{}) string {
	switch msg := message.(type) {
	case error:
		return msg.Error()
	case string:
		return msg
	default:
		return fmt.Sprintf("message: %v has unknown type: %v", message, msg)
	}
}

func Debug(message interface{}, args ...interface{}) {
	msg(DebugLevel, message, args...)
}

func Info(message string, args ...interface{}) {
	msg(InfoLevel, message, args...)
}

func Warn(message string, args ...interface{}) {
	msg(WarnLevel, message, args...)
}

func Error(message interface{}, args ...interface{}) {
	msg(ErrorLevel, message, args...)
}

func Fatal(message interface{}, args ...interface{}) {
	msg(FatalLevel, message, args...)
	os.Exit(1)
}

func logDebug(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Debug().Msg(message)
	} else {
		logger.Debug().Msgf(message, args...)
	}
}

func logInfo(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Info().Msg(message)
	} else {
		logger.Info().Msgf(message, args...)
	}
}

func logWarning(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Warn().Msg(message)
	} else {
		logger.Warn().Msgf(message, args...)
	}
}

func logError(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Error().Msg(message)
	} else {
		logger.Error().Msgf(message, args...)
	}
}

func logfatal(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Fatal().Msg(message)
	} else {
		logger.Fatal().Msgf(message, args...)
	}
}

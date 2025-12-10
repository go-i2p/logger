package logger

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log      *Logger
	once     sync.Once
	failFast string
)

// Logger wraps logrus.Logger and adds the ability to make all warnings fatal
type Logger struct {
	*logrus.Logger
}

// Entry wraps logrus.Entry and enables it to use our Logger
type Entry struct {
	Logger
	entry *logrus.Entry
}

type Fields logrus.Fields

type Level logrus.Level

type TextFormatter logrus.TextFormatter

// Warn wraps logrus.Warn and logs a fatal error if failFast is set
func (l *Logger) Warn(args ...interface{}) {
	warnFatal(args...)
	l.Logger.Warn(args...)
}

// Warnf wraps logrus.Warnf and logs a fatal error if failFast is set
func (l *Logger) Warnf(format string, args ...interface{}) {
	warnFatalf(format, args...)
	l.Logger.Warnf(format, args...)
}

// Error wraps logrus.Error and logs a fatal error if failFast is set
func (l *Logger) Error(args ...interface{}) {
	warnFatal(args...)
	l.Logger.Error(args...)
}

// Errorf wraps logrus.Errorf and logs a fatal error if failFast is set
func (l *Logger) Errorf(format string, args ...interface{}) {
	warnFatalf(format, args...)
	l.Logger.Errorf(format, args...)
}

// WithField wraps logrus.WithField and returns an Entry
func (l *Logger) WithField(key string, value interface{}) *Entry {
	entry := l.Logger.WithField(key, value)
	return &Entry{*l, entry}
}

// WithFields wraps logrus.WithFields and returns an Entry
func (l *Logger) WithFields(fields Fields) *Entry {
	entry := l.Logger.WithFields(logrus.Fields(fields))
	return &Entry{*l, entry}
}

// WithError wraps logrus.WithError and returns an Entry
func (l *Logger) WithError(err error) *Entry {
	entry := l.Logger.WithError(err)
	return &Entry{*l, entry}
}

func (l *Logger) SetLevel(level Level) {
	l.Logger.SetLevel(logrus.Level(level))
}

func (l *Logger) GetLevel() Level {
	return Level(l.Logger.GetLevel())
}

func (l *Logger) SetFormatter(formatter *TextFormatter) {
	l.Logger.SetFormatter((*logrus.TextFormatter)(formatter))
}

func warnFatal(args ...interface{}) {
	if failFast != "" {
		if log != nil {
			log.Fatal(args...)
		} else {
			// Fallback to os.Exit if log is not initialized
			panic("Logger not initialized but fast-fail mode enabled")
		}
	}
}

func warnFatalf(format string, args ...interface{}) {
	if failFast != "" {
		if log != nil {
			log.Fatalf(format, args...)
		} else {
			// Fallback to os.Exit if log is not initialized
			panic("Logger not initialized but fast-fail mode enabled")
		}
	}
}

func warnFail() {
	if failFast != "" {
		log.Error("FATAL ERROR")
	}
}

// InitializeGoI2PLogger sets up all the necessary logging
func InitializeGoI2PLogger() {
	once.Do(func() {
		log = &Logger{}
		log.Logger = logrus.New()
		fmtter := &TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            false,
			DisableColors:          false,
			DisableQuote:           false,
			DisableTimestamp:       false,
			DisableSorting:         false,
			DisableLevelTruncation: false,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:        "time",
				logrus.FieldKeyLevel:       "level",
				logrus.FieldKeyMsg:         "msg",
				logrus.FieldKeyLogrusError: "logrus_error",
				logrus.FieldKeyFunc:        "func",
				logrus.FieldKeyFile:        "file",
			},
		}
		// Configure TextFormatter to include all structured fields
		log.SetFormatter(fmtter)

		// We do not want to log by default
		log.SetOutput(io.Discard)
		log.SetLevel(PanicLevel)
		// Check if DEBUG_I2P is set
		if logLevel := os.Getenv("DEBUG_I2P"); logLevel != "" {
			failFast = os.Getenv("WARNFAIL_I2P")
			if failFast != "" && logLevel == "" {
				logLevel = "debug"
			}
			log.SetOutput(os.Stdout)
			switch strings.ToLower(logLevel) {
			case "debug":
				log.SetLevel(DebugLevel)
			case "warn":
				log.SetLevel(WarnLevel)
			case "error":
				log.SetLevel(ErrorLevel)
			default:
				log.SetLevel(DebugLevel)
			}
			log.WithField("level", log.GetLevel()).Debug("Logging enabled.")
		}
	})
}

// GetGoI2PLogger returns the initialized Logger
func GetGoI2PLogger() *Logger {
	if log == nil {
		InitializeGoI2PLogger()
	}
	return log
}

func init() {
	InitializeGoI2PLogger()
}

var (
	PanicLevel Level = Level(logrus.PanicLevel)
	FatalLevel Level = Level(logrus.FatalLevel)
	ErrorLevel Level = Level(logrus.ErrorLevel)
	WarnLevel  Level = Level(logrus.WarnLevel)
	InfoLevel  Level = Level(logrus.InfoLevel)
	DebugLevel Level = Level(logrus.DebugLevel)
	TraceLevel Level = Level(logrus.TraceLevel)
)

func New() *Logger {
	l := &Logger{}
	l.Logger = logrus.New()
	return l
}

// WithFields creates a new Entry with the specified fields
// INEFFICIENT - use Logger.WithFields instead
func WithFields(fields Fields) *Entry {
	l := GetGoI2PLogger()
	entry := l.Logger.WithFields(logrus.Fields(fields))
	return &Entry{*l, entry}
}

// WithField creates a new Entry with the specified field
// INEFFICIENT - use Logger.WithField instead
func WithField(key string, value interface{}) *Entry {
	l := GetGoI2PLogger()
	entry := l.Logger.WithField(key, value)
	return &Entry{*l, entry}
}

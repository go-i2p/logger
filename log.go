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

// Entry methods that delegate to the underlying logrus.Entry

// Debug logs a debug message with the entry's fields
func (e *Entry) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}

// Debugf logs a formatted debug message with the entry's fields
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

// Info logs an info message with the entry's fields
func (e *Entry) Info(args ...interface{}) {
	e.entry.Info(args...)
}

// Infof logs a formatted info message with the entry's fields
func (e *Entry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Warn logs a warning message with the entry's fields and checks failFast
func (e *Entry) Warn(args ...interface{}) {
	warnFatal(args...)
	e.entry.Warn(args...)
}

// Warnf logs a formatted warning message with the entry's fields and checks failFast
func (e *Entry) Warnf(format string, args ...interface{}) {
	warnFatalf(format, args...)
	e.entry.Warnf(format, args...)
}

// Error logs an error message with the entry's fields and checks failFast
func (e *Entry) Error(args ...interface{}) {
	warnFatal(args...)
	e.entry.Error(args...)
}

// Errorf logs a formatted error message with the entry's fields and checks failFast
func (e *Entry) Errorf(format string, args ...interface{}) {
	warnFatalf(format, args...)
	e.entry.Errorf(format, args...)
}

// Fatal logs a fatal message with the entry's fields and exits
func (e *Entry) Fatal(args ...interface{}) {
	e.entry.Fatal(args...)
}

// Fatalf logs a formatted fatal message with the entry's fields and exits
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

// Panic logs a panic message with the entry's fields and panics
func (e *Entry) Panic(args ...interface{}) {
	e.entry.Panic(args...)
}

// Panicf logs a formatted panic message with the entry's fields and panics
func (e *Entry) Panicf(format string, args ...interface{}) {
	e.entry.Panicf(format, args...)
}

// WithField adds a field to the entry and returns a new Entry
func (e *Entry) WithField(key string, value interface{}) *Entry {
	newEntry := e.entry.WithField(key, value)
	return &Entry{e.Logger, newEntry}
}

// WithFields adds multiple fields to the entry and returns a new Entry
func (e *Entry) WithFields(fields Fields) *Entry {
	newEntry := e.entry.WithFields(logrus.Fields(fields))
	return &Entry{e.Logger, newEntry}
}

// WithError adds an error field to the entry and returns a new Entry
func (e *Entry) WithError(err error) *Entry {
	newEntry := e.entry.WithError(err)
	return &Entry{e.Logger, newEntry}
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
			DisableQuote:           true,
			DisableTimestamp:       false,
			DisableSorting:         false,
			DisableLevelTruncation: false,
			QuoteEmptyFields:       true,
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

package logger

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel int8 = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var (
	conf = zap.NewProductionEncoderConfig()
)

type logger struct {
	core           *zap.SugaredLogger
	enable         bool
	name           string
	level          int8
	encoder        zapcore.Encoder
	outputConsole  bool
	outputFile     bool
	outputFilePath string
	writers        []io.Writer
}

var log = &logger{
	name:    "app",
	encoder: zapcore.NewJSONEncoder(conf),
}

func init() {
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
}

func (l *logger) init() {
	l.initOutputConsole()
	l.initOutputFile()

	cores := make([]zapcore.Core, len(l.writers))
	for i, w := range l.writers {
		cores[i] = zapcore.NewCore(l.encoder, zapcore.Lock(zapcore.AddSync(w)), zapcore.Level(l.level))
	}
	l.core = zap.New(zapcore.NewTee(cores...)).Sugar()
}

func (l *logger) initOutputFile() {
	if l.outputFile {
		if _, err := os.Stat(l.outputFilePath); os.IsNotExist(err) {
			os.Mkdir(l.outputFilePath, os.ModePerm)
		}
		f, _ := os.OpenFile(fmt.Sprintf("%s/%s.log", l.outputFilePath, l.name), os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
		l.writers = append(l.writers, f)
	}
}

func (l *logger) initOutputConsole() {
	if l.outputConsole {
		l.writers = append(l.writers, os.Stdout)
	}
}

// New returns a logger by specify a name
func New(name string) *logger {
	return &logger{
		name:           name,
		outputFilePath: "./logs",
	}
}

// Enable makes logger enable to use.
func Enable() *logger {
	log.enable = true
	log.init()
	return log
}

// Enable makes logger enable to use.
func (l *logger) Enable() *logger {
	l.enable = true
	l.init()
	return l
}

// Level sets the level of logger.
// The default level is `Info`.
func Level(level int8) *logger {
	log.level = level
	return log
}

// Level sets the level of logger.
// The default level is `Info`.
func (l *logger) Level(level int8) *logger {
	l.level = level
	return l
}

// ConsoleEncoding EncodeConsole sets the encoding `console` of logger.
// The default encoding is `json`.
func ConsoleEncoding() *logger {
	return log.EncodeConsole()
}

// EncodeConsole sets the encoding `console` of logger.
// The default encoding is `json`.
func (l *logger) EncodeConsole() *logger {
	l.encoder = zapcore.NewConsoleEncoder(conf)
	return l
}

// WithFile sets the output file of logger.
func WithFile(path string) *logger {
	return log.WithFile(path)
}

// WithFile sets the output file of logger.
func (l *logger) WithFile(path string) *logger {
	l.outputFile = true
	l.outputFilePath = path
	return l
}

func WithConsole() *logger {
	return log.WithConsole()
}

func (l *logger) WithConsole() *logger {
	l.outputConsole = true
	return l
}

func WithWriter(w io.Writer) *logger {
	return log.WithWriter(w)
}

func (l *logger) WithWriter(w io.Writer) *logger {
	l.writers = append(l.writers, w)
	return l
}

func Debug(msg string, keyAndValues ...interface{}) {
	log.Debug(msg, keyAndValues...)
}

func (l *logger) Debug(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Debugw(msg, keyAndValues...)
	}
}

func Info(msg string, keyAndValues ...interface{}) {
	log.Info(msg, keyAndValues...)
}

func (l *logger) Info(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Infow(msg, keyAndValues...)
	}
}

func Warn(msg string, keyAndValues ...interface{}) {
	log.Warn(msg, keyAndValues...)
}

func (l *logger) Warn(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Warnw(msg, keyAndValues...)
	}
}

func Error(msg string, keyAndValues ...interface{}) {
	log.Error(msg, keyAndValues)
}

func (l *logger) Error(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Errorw(msg, keyAndValues...)
	}
}

func Panic(msg string, keyAndValues ...interface{}) {
	log.Panic(msg, keyAndValues...)
}

func (l *logger) Panic(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Panicw(msg, keyAndValues...)
	}
}

func Fatal(msg string, keyAndValues ...interface{}) {
	log.Fatal(msg, keyAndValues...)
}

func (l *logger) Fatal(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Fatalw(msg, keyAndValues...)
	}
}

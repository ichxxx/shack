package shack

import (
	"fmt"
	"os"
	"strings"

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

type logger struct {
	core        *zap.SugaredLogger
	enable      bool
	name        string
	level       int8
	encoding    string
	outputPaths []string
	development bool
}

var Log = &logger{}


func(l *logger) init() {
	var writeSyncer zapcore.WriteSyncer
	for _, path := range l.outputPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
		f, _ := os.OpenFile(fmt.Sprintf("%s/%s.log", path, l.name), os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
		writeSyncer = zapcore.AddSync(f)
	}

	var conf zapcore.EncoderConfig
	if l.development {
		conf = zap.NewDevelopmentEncoderConfig()
	} else {
		conf = zap.NewProductionEncoderConfig()
	}
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	var encoder zapcore.Encoder
	switch strings.ToLower(l.encoding) {
	case "console":
		encoder = zapcore.NewConsoleEncoder(conf)
	default:
		encoder = zapcore.NewJSONEncoder(conf)
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(l.level))
	l.core = zap.New(core).Sugar()
}


// Enable makes logger enable to use.
func(l *logger) Enable() {
	Log.enable = true
	Log.init()
}


// Level sets the level of logger.
// The default level is `Info`.
func(l *logger) Level(level int8) *logger {
	l.level = level
	return l
}


// Level sets the encoding ( `json` or `console` ) of logger.
// The default encoding is `json`.
func(l *logger) Encoding(encoding string) *logger {
	l.encoding = encoding
	return l
}


// Output sets the output paths of logger.
// The default output path is `./logs`.
func(l *logger) Output(paths ...string) *logger {
	l.outputPaths = append(l.outputPaths, paths...)
	return l
}


// Dev enable the development mode of logger.
func(l *logger) Dev() *logger {
	l.development = true
	return l
}


func(l *logger) Debug(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Debugw(msg, keyAndValues...)
	}
}


func(l *logger) Info(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Infow(msg, keyAndValues...)
	}
}


func(l *logger) Warn(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Warnw(msg, keyAndValues...)
	}
}


func(l *logger) Error(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Errorw(msg, keyAndValues...)
	}
}


func(l *logger) Panic(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Panicw(msg, keyAndValues...)
	}
}


func(l *logger) Fatal(msg string, keyAndValues ...interface{}) {
	if l.enable {
		l.core.Fatalw(msg, keyAndValues...)
	}
}

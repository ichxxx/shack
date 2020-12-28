package middleware

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ichxxx/shack"
)


func initAccessLogger(path string) *zap.Logger {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	f, _ := os.OpenFile(path + "/access.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	writeSyncer := zapcore.AddSync(f)
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoder := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	return zap.New(core)
}


func AccessLog(path ...string) shack.HandlerFunc {
	var accessLogger *zap.Logger
	if len(path) > 0 {
		accessLogger = initAccessLogger(strings.Join(path, ""))
	} else {
		accessLogger = initAccessLogger("./logs")
	}

	return func(ctx *shack.Context) {
		start := time.Now()
		duration := time.Since(start)
		durationMs := duration.Nanoseconds() / (1000 * 1000)

		accessLogger.Info("",
			zap.String("uri", ctx.Path),
			zap.String("request_method", ctx.Method),
			zap.String("query_string", ctx.Request.URL.RawQuery),
			zap.Int("status", ctx.StatusCode),
			zap.Int64("response_time", durationMs),
			zap.String("remote_address", ctx.Request.RemoteAddr),
			zap.String("protocol", ctx.Request.Proto),
			zap.String("server_name", ctx.Request.URL.Host),
			zap.Int64("bytes_received", ctx.Request.ContentLength),
		)

		ctx.Next()
	}
}

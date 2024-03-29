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
	f, _ := os.OpenFile(path+"/access.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	writeSyncer := zapcore.AddSync(f)
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoder := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	return zap.New(core)
}

func AccessLog(path ...string) shack.Handler {
	var accessLogger *zap.Logger
	if len(path) > 0 {
		accessLogger = initAccessLogger(strings.Join(path, ""))
	} else {
		accessLogger = initAccessLogger("./logs")
	}

	return func(ctx *shack.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		durationMs := float64(duration.Nanoseconds()) / (1000 * 1000)

		accessLogger.Info("",
			zap.Float64("response_ms", durationMs),
			zap.String("path", ctx.Request.Path()),
			zap.String("method", ctx.Request.Method()),
			zap.String("query", ctx.Request.RawQuery()),
			zap.Int("code", ctx.Response.StatusCode),
			zap.String("remote_address", ctx.Request.RemoteAddr),
			zap.String("server_name", ctx.Request.Host),
		)
	}
}

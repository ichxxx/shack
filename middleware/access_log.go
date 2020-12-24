package middleware

import (
	"time"

	"go.uber.org/zap"

	"shack"
)

var AccessLogger, _ = zap.NewProduction()


func AccessLog(ctx *shack.Context) {
	start := time.Now()
	defer AccessLogger.Sync()

	duration := time.Since(start)
	durationMs := duration.Nanoseconds() / (1000 * 1000)

	AccessLogger.Info("",
		zap.String("log_type", "access"),
		zap.String("remote_address", ctx.Request.RemoteAddr),
		zap.Int64("response_time", durationMs),
		zap.String("protocol", ctx.Request.Proto),
		zap.String("request_method", ctx.Method),
		zap.String("query_string", ctx.Request.URL.RawQuery),
		zap.Int("status", ctx.StatusCode),
		zap.String("uri", ctx.Path),
		zap.String("server_name", ctx.Request.URL.Host),
		zap.Int64("bytes_received", ctx.Request.ContentLength),
	)

	ctx.Next()
}

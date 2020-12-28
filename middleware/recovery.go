package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/ichxxx/shack"
)


func Recovery() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				ctx.Status(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}


func trace(message string) string {
	var pcs [32]uintptr
	// skip first 3 caller
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n%s:%d", file, line))
	}
	return str.String()
}

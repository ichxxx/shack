package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/ichxxx/shack"
)


func Recovery() shack.HandlerFunc {
	return func(ctx *shack.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s\n\n", trace(fmt.Sprintf("%s", err)))
				ctx.HttpStatus(http.StatusInternalServerError)
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
	str.WriteString("\nError:\n\t")
	str.WriteString(message)
	str.WriteString("\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString("\n\t")
		str.WriteString(file)
		str.WriteString(":")
		str.WriteString(strconv.Itoa(line))
	}
	return str.String()
}

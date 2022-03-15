package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// ChiZapLoggerMiddleware is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func ChiZapLoggerMiddleware(l *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqStartTime := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				l.Info("Request",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("params", r.URL.RawQuery),
					zap.Int("status", ww.Status()),
					zap.String("size", fmt.Sprintf("%dB", ww.BytesWritten())),
					zap.Duration("lat", time.Since(reqStartTime)),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

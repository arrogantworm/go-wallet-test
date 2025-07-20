package middlewares

import (
	"log"
	"net/http"
	"time"
)

type logger struct {
	http.ResponseWriter
	statusCode int
}

func (l *logger) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := &logger{
			ResponseWriter: w,
			statusCode:     200,
		}
		next.ServeHTTP(l, r)
		log.Printf("LOG: %s %s %d %s\n", r.Method, r.URL.Path, l.statusCode, time.Since(start))
	})
}

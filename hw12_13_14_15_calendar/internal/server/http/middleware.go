package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(logger *logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		msg := fmt.Sprintf("%s %s %s %s %s %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			r.Proto,
			time.Since(start).String(),
			r.UserAgent())
		logger.Info(msg)
	}
}

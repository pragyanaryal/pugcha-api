package middleware

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/log"
	"net/http"
	"time"
)

// LogEachRequest ...
func LogEachRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(rw, r)
		log.Info().
			Str("Method", r.Method).
			Str("Host", r.Host).
			Str("Request URI", r.RequestURI).
			Dur("Latency", time.Since(startTime)).
			Msgf("time", time.Since(startTime))
	})
}

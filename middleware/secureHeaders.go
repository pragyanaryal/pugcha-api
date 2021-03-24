package middleware

import "net/http"

// SecureHeaders ...
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// Set the header of each request to the server
		rw.Header().Set("X-XSS-Protection", "1; mode=block")
		rw.Header().Set("X-Frame-Options", "deny")

		// log the setting of secure header

		next.ServeHTTP(rw, r)
	})
}

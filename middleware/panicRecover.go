package middleware

import "net/http"

// PanicRecover ...
func PanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		defer func() {
			// use the built in recover function to check if there has been panic or not.
			if err := recover(); err != nil {
				// Set the connection close header on the response
				rw.Header().Set("Connection", "close")

				// Send the ServerError response back to client

			}
		}()
		next.ServeHTTP(rw, r)
	})
}

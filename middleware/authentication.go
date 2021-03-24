package middleware

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
	"strings"
)

type AuthUserTemp struct {
	UserId        uuid.UUID
	UserType      string
	LoggedFrom    string
	Authenticated bool
}

// CheckAuthentication ...
func CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		auth := AuthUserTemp{}
		var req *http.Request

		// Get access token
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) > 1 {
			claims, errs := serveUtils.DecodeJWT(splitToken[1], "access")
			if errs == nil {
				auth.UserId = uuid.MustParse((*claims)["userId"].(string))
				auth.UserType = (*claims)["type"].(string)
				auth.LoggedFrom = (*claims)["loggedFrom"].(string)
				auth.Authenticated = true

				req = getNewRequest(r, auth)
			} else {
				req = getNewRequest(r, auth)
			}
		} else {
			req = getNewRequest(r, auth)
		}
		next.ServeHTTP(rw, req)
	})
}

func getNewRequest(r *http.Request, auth AuthUserTemp) *http.Request {
	//auth.Authenticated = true
	//auth.UserType = "super"
	ctx := context.WithValue(r.Context(), AuthUserTemp{}, auth)
	return r.WithContext(ctx)
}

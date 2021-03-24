package serveUtils

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"net/http"
	"time"
)

func SetCookie(rw http.ResponseWriter, token *json_serializer.TokenSerializer) {

	access_cookie := http.Cookie{
		Name:     "pugcha_token",
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 2),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	}

	refresh_cookie := http.Cookie{
		Name:     "pugcha_refresh",
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 10),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	}

	http.SetCookie(rw, &access_cookie)
	http.SetCookie(rw, &refresh_cookie)
}

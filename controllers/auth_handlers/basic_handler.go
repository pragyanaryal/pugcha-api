package auth_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/use_cases"
	"net/http"
)

func (auth *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	loginCredential := json_serializer.LoginSerializer{}

	err := json.NewDecoder(r.Body).Decode(&loginCredential)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	err, _, token := use_cases.LoginUserOrchestrator(rw, [16]byte{}, "", "normal", &loginCredential)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}
	//user.AccessToken = token.AccessToken
	//user.RefreshToken = token.RefreshToken
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(rw).Encode(token)
}

func (auth *AuthHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {

		accessCookie := http.Cookie{
			Name:   "pugcha_token",
			Path:   "/",
			Value:  "",
			MaxAge: -1,
		}
		refreshCookie := http.Cookie{
			Name:   "pugcha_refresh",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(rw, &accessCookie)
		http.SetCookie(rw, &refreshCookie)

	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
}

func (auth *AuthHandler) ChangePassword(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {

	}
	http.Error(rw, "Not Authenticated", http.StatusUnauthorized)
}

func (auth *AuthHandler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
}

func (auth *AuthHandler) VerifyAccount(rw http.ResponseWriter, r *http.Request) {
	param, err := getParams(r, "id", "token")
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
	}

	uID, err := uuid.Parse(param[0])
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	password := json_serializer.PasswordSerializer{}

	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	user, err := use_cases.VerifyAccountOrchestrator(uID, param[1], &password)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(rw).Encode(user)
	return
}

func (auth *AuthHandler) ConfirmPassword(rw http.ResponseWriter, r *http.Request) {
	param, err := getParams(r, "id", "token")
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
	}

	uID, err := uuid.Parse(param[0])
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	password := json_serializer.PasswordSerializer{}

	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	_, _, err = use_cases.ResetPasswordOrchestrator(uID, param[1], &password, "reset")
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	return
}

func (auth *AuthHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	token, err := serveUtils.RefreshToken(r)

	if err != nil {
		erros.JSONError(rw, err, http.StatusNotAcceptable)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(rw).Encode(token)
}

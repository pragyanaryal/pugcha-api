package auth_handlers

import (
	"encoding/json"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/use_cases"
	"net/http"
)

func (auth *AuthHandler) GoogleLogin(rw http.ResponseWriter, r *http.Request) {
	serveUsers.GoogleService.GoogleLogin(rw, r)
}

func (auth *AuthHandler) GoogleCallback(rw http.ResponseWriter, r *http.Request) {
	user, err := use_cases.GoogleCallBackOrchestrator(r)
	if err != nil {
		erros.JSONError(rw, err, http.StatusInternalServerError)
		return
	}

	err, _, _ = use_cases.LoginUserOrchestrator(rw, user.UserId, "user", "google", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(rw).Encode(user)
	if e != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

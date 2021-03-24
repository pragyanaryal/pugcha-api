package auth_handlers

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/service/use_cases"
	"net/http"
)

func (auth *AuthHandler) FacebookLogin(rw http.ResponseWriter, r *http.Request) {

}

func (auth *AuthHandler) FacebookCallback(rw http.ResponseWriter, r *http.Request) {
	use_cases.FacebookCallBackOrchestrator()
}

package auth_handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type AuthHandler struct{}

func NewAuth() *AuthHandler { return &AuthHandler{} }

func getParams(r *http.Request, args ...interface{}) ([]string, error) {
	var ar []string

	for _, arg := range args {
		vars := mux.Vars(r)

		temp, ok := vars[arg.(string)]
		if !ok {
			return nil, errors.New("not found")
		}

		ar = append(ar, temp)
	}

	return ar, nil
}

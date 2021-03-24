package business_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveBusiness"
	"net/http"
)

func (c *BusinessHandler) CreateBusinessProfile(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {

	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) GetBusinessesProfile(rw http.ResponseWriter, r *http.Request) {
	service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	profile, err := service.ListBusinessProfiles()
	if err != nil {
		erros.JSONError(rw, err, http.StatusNotFound)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(profile)
	return
}

func (c *BusinessHandler) GetBusinessProfile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		erros.JSONError(rw, errors.New("unable to convert id"), http.StatusBadRequest)
		return
	}

	service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	profile, err := service.FindProfileById(uID)
	if err != nil {
		erros.JSONError(rw, err, http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(profile)
	return
}

func (c *BusinessHandler) UpdateBusinessProfile(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {

	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return

}

func (c *BusinessHandler) DeleteBusinessProfile(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		if userInfo.UserType == "admin" {

		}
		erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
		return

	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return

}

package user_handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"net/http"
)

// CreateUsers ...
func (u *UsersHandler) CreateUserProfile(rw http.ResponseWriter, r *http.Request) {
}

func (u *UsersHandler) GetUserProfiles(rw http.ResponseWriter, r *http.Request) {

	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	fmt.Println(userInfo, "userinfo")

	if userInfo.UserId == uuid.Nil {
		erros.JSONError(rw, errors.New("no token sent"), http.StatusForbidden)
		return
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	user, err := service.GetUserById(userInfo.UserId)
	if err != nil {
		erros.JSONError(rw, err, http.StatusForbidden)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(user)

}

func (u *UsersHandler) GetSingleProfile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		erros.JSONError(rw, errors.New("unable to convert id"), http.StatusBadRequest)
		return
	}
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	if userInfo.UserId == uID || userInfo.UserType == "admin" || userInfo.UserType == "super" {
		fmt.Println("here?", userInfo.UserId, userInfo.UserType)
		var (
			userRepo    = repositories.UserRepo
			profileRepo = repositories.UserProfileRepo
			googleRepo  = repositories.GoogleRepo
			fbRepo      = repositories.FacebookRepo
		)

		service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

		profile, err := service.GetProfileById(uID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		err = json.NewEncoder(rw).Encode(profile)
		if err != nil {
			erros.JSONError(rw, err, http.StatusInternalServerError)
			return
		}
		return
	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (u *UsersHandler) UpdateUserProfile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	if userInfo.UserId == uID || userInfo.UserType == "admin" || userInfo.UserType == "super" {

	}
}

func (u *UsersHandler) DeleteUserProfile(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	if userInfo.UserType == "admin" || userInfo.UserType == "super" {

	}
	http.Error(rw, "Not enough authority", http.StatusUnauthorized)
}

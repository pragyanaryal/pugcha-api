package user_handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/use_cases"
	"net/http"
	"strings"
)

func (u *UsersHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	userType := r.Header.Get("user_type")

	if userType == "normal" {
		normal := json_serializer.CreateNormalUser{}

		e := json.NewDecoder(r.Body)
		err := e.Decode(&normal)
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		user, err := use_cases.HandleNormalUserOrchestrator(&normal, rw)
		if err != nil {
			erros.JSONError(rw, err, http.StatusInternalServerError)
			return
		}
		err, _, token := use_cases.LoginUserOrchestrator(rw, user.ID, "user", normal.SocialPlatform, nil)
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(rw).Encode(token)
		return
	}

	if userInfo.UserType == "super" || userType == "staff" {
		staff := json_serializer.CreateAdminUserRequest{}

		e := json.NewDecoder(r.Body)
		err := e.Decode(&staff)
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		user, err := use_cases.CreateStaffUserOrchestrator(&staff)
		if err != nil {

			if err.Error() == "an user with email already exists" {
				erros.JSONError(rw, err, http.StatusConflict)
				return
			}
			if err.Error() == "validation error" {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}
			erros.JSONError(rw, err, http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(rw).Encode(user)
		return
	}
	erros.JSONError(rw, errors.New("unauthorized"), http.StatusUnauthorized)
	return
}

func (u *UsersHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	user, err := use_cases.ListUsers(userInfo, r)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(user)
	return
}

func (u *UsersHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	user, err := service.GetUserById(uID)
	if err != nil {
		erros.JSONError(rw, err, http.StatusForbidden)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(user)
}

func (u *UsersHandler) UpdateUser(rw http.ResponseWriter, r *http.Request) {

}

func (u *UsersHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	update := false

	if userInfo.Authenticated == true {
		update = true
	} else if r.Header.Get("check") == "verify-update" {
		vars := mux.Vars(r)
		uID, err := uuid.Parse(vars["id"])

		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		if r.Header.Get("user") == uID.String() {
			var (
				userRepo    = repositories.UserRepo
				profileRepo = repositories.UserProfileRepo
				googleRepo  = repositories.GoogleRepo
				fbRepo      = repositories.FacebookRepo
			)

			service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

			user, err := service.GetUserById(uID)
			if err != nil {
				erros.JSONError(rw, err, http.StatusForbidden)
				return
			}
			if user.Profile == nil {
				userInfo.UserId = user.ID
				update = true
			}
		}
	}

	if update {
		vars := mux.Vars(r)
		uID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		if userInfo.UserId == uID || userInfo.UserType == "admin" || userInfo.UserType == "super" {

			user := json_serializer.UpdateAbleUser{}

			var decoder *form.Decoder
			decoder = form.NewDecoder()

			content := strings.Split(r.Header.Get("Content-Type"), ";")[0]

			if content == "multipart/form-data" {
				_ = r.ParseMultipartForm(32 << 20)
				err := decoder.Decode(&user, r.PostForm)

				if err != nil {
					erros.JSONError(rw, err, http.StatusBadRequest)
					return
				}
			} else if content == "application/x-www-form-urlencoded" {
				_ = r.ParseForm()
				err := decoder.Decode(&user, r.PostForm)

				if err != nil {
					erros.JSONError(rw, err, http.StatusBadRequest)
					return
				}
			}

			if user.Status != nil || user.Type != nil {
				if userInfo.UserType == "super" || userInfo.UserType == "admin" && userInfo.UserId == uID {
					users, err := use_cases.UpdateUser(&user, uID, content, r)
					if err != nil {
						erros.JSONError(rw, err, http.StatusBadRequest)
						return
					}
					rw.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(rw).Encode(users)
					return
				}
				erros.JSONError(rw, errors.New("user is not authorized"), http.StatusUnauthorized)
				return
			}
			users, err := use_cases.UpdateUser(&user, uID, content, r)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(users)
			return
		}
	}
	erros.JSONError(rw, errors.New("user is not authenticated"), http.StatusUnauthorized)
	return
}

func (u *UsersHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		_, _ = rw.Write([]byte("Delete Users"))
	}
	http.Error(rw, "Not Authenticated", http.StatusUnauthorized)

}

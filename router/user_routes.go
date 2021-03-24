package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/controllers"
	"net/http"
)

var user = controllers.NewUsers()

func SetupUserRoutes(r *mux.Router) {

	r.HandleFunc("/user", user.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/user", user.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", user.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", user.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/user/{id}", user.PatchUser).Methods(http.MethodPatch)
	r.HandleFunc("/user/{id}", user.DeleteUser).Methods(http.MethodDelete)

	// profile
	r.HandleFunc("/user-profile", user.GetUserProfiles).Methods(http.MethodGet)
	r.HandleFunc("/user-profile", user.CreateProfile).Methods(http.MethodPost)
}

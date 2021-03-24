package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/controllers"
	"net/http"
)

var auth = controllers.NewAuth()

func SetupAuthRoutes(r *mux.Router) {
	// Regular Routes
	r.HandleFunc("/auth/login", auth.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodPost)
	r.HandleFunc("/auth/changePassword", auth.ChangePassword).Methods(http.MethodPost)
	r.HandleFunc("/auth/resetPassword", auth.ResetPassword).Methods(http.MethodGet)
	r.HandleFunc("/auth/verify/{id}/{token}", auth.VerifyAccount).Methods(http.MethodPost)
	r.HandleFunc("/auth/confirmPassword/{id}/{token}", auth.ConfirmPassword).Methods(http.MethodPost)

	// Token Refresh
	r.HandleFunc("/auth/refreshToken", auth.RefreshToken).Methods(http.MethodGet)
}

package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"net/http"
)

var (
	r = mux.NewRouter()
)

func New() http.Handler {
	SetupFrontRoutes(r, config.Configuration.FrontHost)
	SetupAdminRoutes(r, config.Configuration.AdminHost)

	s := r.Host(config.Configuration.Host).Schemes("https").Subrouter()

	SetupAuthRoutes(s)
	SetupUserRoutes(s)
	SetupCategoryRoutes(s)
	SetupBusinessRoutes(s)
	SetupAddressRoutes(s)
	SetupSearchRoutes(s)
	SetupFileRoutes(s, config.Configuration.File, config.Configuration.Host)

	s.Use(middleware.PanicRecover)
	s.Use(middleware.LogEachRequest)
	s.Use(middleware.SecureHeaders)
	s.Use(middleware.CheckAuthentication)
	s.Use(middleware.GzipMiddleware)

	return r
}

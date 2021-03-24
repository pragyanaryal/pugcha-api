package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/controllers"
	"net/http"
)

var search = controllers.NewSearch()

func SetupSearchRoutes(r *mux.Router) {
	r.HandleFunc("/search/{term}", search.Search).Methods(http.MethodGet)
}


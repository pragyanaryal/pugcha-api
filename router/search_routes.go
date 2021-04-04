package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/handlers"
	"net/http"
)

var search = handlers.NewSearch()

func SetupSearchRoutes(r *mux.Router) {
	r.HandleFunc("/search/{term}", search.Search).Methods(http.MethodGet)
}


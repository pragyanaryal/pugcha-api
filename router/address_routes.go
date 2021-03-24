package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/controllers"
	"net/http"
)

var address = controllers.NewAddress()


func SetupAddressRoutes(r *mux.Router) {
	r.HandleFunc("/business/{id}/address", address.AddAddress).Methods(http.MethodPost)
	r.HandleFunc("/business/{id}/address/", address.UpdateAddress).Methods(http.MethodPut)
	r.HandleFunc("/business/{id}/address/{pid}", address.PatchAddress).Methods(http.MethodPatch)
	r.HandleFunc("/business/{id}/address/{pid}", address.DeleteAddress).Methods(http.MethodDelete)
}


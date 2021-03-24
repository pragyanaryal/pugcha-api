package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/controllers"
	"net/http"
)

var business = controllers.NewBusinessHandler()

func SetupBusinessRoutes(r *mux.Router) {

	r.HandleFunc("/business", business.CreateBusiness).Methods(http.MethodPost).Name("create business")
	r.HandleFunc("/business", business.GetBusinesses).Methods(http.MethodGet)
	r.HandleFunc("/business/{id}", business.GetBusiness).Methods(http.MethodGet)
	r.HandleFunc("/business/{id}", business.PatchBusiness).Methods(http.MethodPatch)
	r.HandleFunc("/business/{id}", business.DeleteBusiness).Methods(http.MethodDelete)
	r.HandleFunc("/business/search", business.SearchBusiness).Methods(http.MethodPost, "GET")

	// profiles
	r.HandleFunc("/business-profile/{id}", business.GetBusinessProfile).Methods(http.MethodGet)

	// pictures
	r.HandleFunc("/business/{id}/picture", business.AddPicture).Methods(http.MethodPost)
	r.HandleFunc("/business/{id}/picture/", business.UpdatePicture).Methods(http.MethodPut)
	r.HandleFunc("/business/{id}/picture/{pid}", business.PatchPicture).Methods(http.MethodPatch)
	r.HandleFunc("/business/{id}/picture/{pid}", business.DeletePicture).Methods(http.MethodDelete)
}

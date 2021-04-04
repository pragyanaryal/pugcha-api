package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/handlers"
	"net/http"
)

var category = handlers.NewCategoriesHandler()

func SetupCategoryRoutes(r *mux.Router) {

	r.HandleFunc("/category", category.CreateCategory).Methods(http.MethodPost)
	r.HandleFunc("/category", category.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/category/{id}", category.GetCategory).Methods(http.MethodGet)
	r.HandleFunc("/category/{id}", category.UpdateCategory).Methods(http.MethodPut)
	r.HandleFunc("/category/{id}", category.DeleteCategory).Methods(http.MethodDelete)
}

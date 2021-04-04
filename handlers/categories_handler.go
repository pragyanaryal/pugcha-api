package handlers

import (
	"encoding/json"
	"net/http"
	
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/category_service"
	"gitlab.com/ProtectIdentity/pugcha-backend/use_cases"
)

type CategoriesHandler struct{}

func NewCategoriesHandler() *CategoriesHandler {
	return &CategoriesHandler{}
}

func (c *CategoriesHandler) CreateCategory(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		if userInfo.UserType == "admin" || userInfo.UserType == "super" {
			category, err := use_cases.CreateCategoriesOrchestrator(r, userInfo.UserId)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(category)
			return
		}
	}
}

func (c *CategoriesHandler) GetCategories(rw http.ResponseWriter, r *http.Request) {
	service := category_service.CategoriesService(repositories.CategoryRepo)

	categories, err := service.GetCategories()
	if err != nil {
		erros.JSONError(rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(categories)
	return
}

func (c *CategoriesHandler) GetCategory(rw http.ResponseWriter, r *http.Request) {
}

func (c *CategoriesHandler) UpdateCategory(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true && userInfo.UserType == "admin" || userInfo.UserType == "super" {

	}
	http.Error(rw, "Not Authenticated", http.StatusUnauthorized)
}

func (c *CategoriesHandler) DeleteCategory(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		if userInfo.UserType == "admin" {

		}
	}
	http.Error(rw, "Not Authenticated", http.StatusUnauthorized)
}

package controllers

import (
	"errors"
	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveBusiness"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/use_cases"
	"net/http"
)

type AddressHandler struct{}

func NewAddress() *AddressHandler { return &AddressHandler{} }


func (add *AddressHandler) AddAddress(rw http.ResponseWriter, r *http.Request){
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindProfileById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {

			newBusiness := json_serializer.UpdateAbleBusiness{}

			var decoder *form.Decoder

			decoder = form.NewDecoder()
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeAddresses, []models.Address{})

			_ = r.ParseForm()
			err := decoder.Decode(&newBusiness, r.PostForm)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}

			if newBusiness.Address != nil {

				profile, err := serveBusiness.PopulateBusinessIdOnUpdate(bID, &newBusiness)
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}

				oldAddress := business.Address

				for _, val := range *profile.Address {
					*oldAddress = append(*oldAddress, val)
				}
				profile.Address = oldAddress

				err = use_cases.UpdateBusiness(r, profile, bID, "nothing")
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusNoContent)
				return
			}
		} else {
			erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
			return
		}

	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (add *AddressHandler) UpdateAddress(rw http.ResponseWriter, r *http.Request){
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindProfileById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {

			newBusiness := json_serializer.UpdateAbleBusiness{}

			var decoder *form.Decoder

			decoder = form.NewDecoder()
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeAddresses, []models.Address{})

			_ = r.ParseForm()
			err := decoder.Decode(&newBusiness, r.PostForm)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}

			if newBusiness.Address != nil {
				profile, err := serveBusiness.PopulateBusinessIdOnUpdate(bID, &newBusiness)

				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}

				var id []uuid.UUID
				for _, val := range *business.Address {
					id = append(id, val.Id)
				}

				err = service.BatchDeleteAddressById(id)
				if err != nil {
					erros.JSONError(rw, errors.New("partial update"), http.StatusInternalServerError)
					return
				}

				err = use_cases.UpdateBusiness(r, profile, bID, "nothing")
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}

				rw.WriteHeader(http.StatusNoContent)
				return
			}
		} else {
			erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
			return
		}

	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (add *AddressHandler) PatchAddress(rw http.ResponseWriter, r *http.Request){

}

func (add *AddressHandler) DeleteAddress(rw http.ResponseWriter, r *http.Request){

}
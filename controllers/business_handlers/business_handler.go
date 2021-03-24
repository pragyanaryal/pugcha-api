package business_handlers

import (
	"encoding/json"
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
	"strings"
	"time"
)

func (c *BusinessHandler) CreateBusiness(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		profile, err := use_cases.CreateBusinessOrchestrator(r)
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(profile)
		return
	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) GetBusinesses(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	business, err := use_cases.ListBusinesses(userInfo, r)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(business)
	return
}

func (c *BusinessHandler) GetBusiness(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	business, err := service.FindById(uID)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(business)
}

func (c *BusinessHandler) UpdateBusiness(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {

	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) PatchBusiness(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}
		service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {
			content := strings.Split(r.Header.Get("Content-Type"), ";")[0]
			business := json_serializer.UpdateAbleBusiness{}
			var decoder *form.Decoder

			decoder = form.NewDecoder()
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeLocation, models.Location{})
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeOwner, []models.Owner{})
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeAddresses, []models.Address{})
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeOpeningHours, []models.OpeningHours{})
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeNullTime, &time.Time{})
			decoder.RegisterCustomTypeFunc(serveBusiness.DecodeUUID, uuid.UUID{})

			if content == "multipart/form-data" {
				r.ParseMultipartForm(32 << 20)
				err := decoder.Decode(&business, r.PostForm)

				if err != nil {
					erros.JSONError(rw, err, http.StatusBadRequest)
					return
				}
			} else if content == "application/x-www-form-urlencoded" {
				r.ParseForm()
				err := decoder.Decode(&business, r.PostForm)

				if err != nil {
					erros.JSONError(rw, err, http.StatusBadRequest)
					return
				}
			}

			if business.Approved != nil || business.UserId != nil || business.ApprovedBy != nil ||
				business.Blocked != nil || business.CategoryId != nil {
				if userInfo.UserType == "super" || userInfo.UserType == "admin" {
					use_cases.UpdateBusiness(r, &business, bID, content)
					return
				}
				erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}
			//businessMap := structomap.New().UseSnakeCase().PickAll().Transform(business)
			//fmt.Println(businessMap, "outside")
			use_cases.UpdateBusiness(r, &business, bID, content)
		}
		erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) DeleteBusiness(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		if userInfo.UserType == "admin" {

		}
		erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
		return
	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

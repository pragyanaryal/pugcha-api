package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	business2 "gitlab.com/ProtectIdentity/pugcha-backend/use_cases/business"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/business_service"
)

type BusinessHandler struct{}

func NewBusinessHandler() *BusinessHandler {
	return &BusinessHandler{}
}

func (c *BusinessHandler) CreateBusiness(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		profile, err := business2.CreateBusinessOrchestrator(r)
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

// Get All Businesses Handler
func (c *BusinessHandler) GetBusinesses(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)

	business, err := business2.ListBusinesses(userInfo, r)
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

	service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	business, err := service.FindById(uID)
	if err != nil {
		erros.JSONError(rw, err, http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(business)
	return
}

func (c *BusinessHandler) GetBusinessProfile(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID, err := uuid.Parse(vars["id"])
	if err != nil {
		erros.JSONError(rw, errors.New("unable to convert id"), http.StatusBadRequest)
		return
	}

	service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	profile, err := service.FindProfileById(uID)
	if err != nil {
		erros.JSONError(rw, err, http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(profile)
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
		service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
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
			decoder.RegisterCustomTypeFunc(business_service.DecodeOwner, []models.Owner{})
			decoder.RegisterCustomTypeFunc(business_service.DecodeLocation, models.Location{})
			decoder.RegisterCustomTypeFunc(business_service.DecodeOpeningHours, []models.OpeningHours{})
			decoder.RegisterCustomTypeFunc(business_service.DecodeNullTime, &time.Time{})
			decoder.RegisterCustomTypeFunc(business_service.DecodeUUID, uuid.UUID{})

			if content == "multipart/form-data" {
				_ = r.ParseMultipartForm(32 << 20)
			} else if content == "application/x-www-form-urlencoded" {
				_ = r.ParseForm()
			} else {
				erros.JSONError(rw, errors.New("invalid request"), http.StatusBadRequest)
				return
			}

			err := decoder.Decode(&business, r.PostForm)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}

			if business.Approved != nil || business.UserId != nil || business.ApprovedBy != nil ||
				business.Blocked != nil || business.CategoryId != nil {
				if userInfo.UserType == "super" || userInfo.UserType == "admin" {
					err = business2.UpdateBusiness(r, &business, bID, content)
					if err != nil {
						erros.JSONError(rw, err, http.StatusBadRequest)
						return
					}
					rw.WriteHeader(http.StatusNoContent)
					return
				}
				erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}
			err = business2.UpdateBusiness(r, &business, bID, content)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}
	erros.JSONError(rw, errors.New("not authenticated"), http.StatusUnauthorized)
	return
}

var upgraded = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *BusinessHandler) SearchBusiness(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgraded.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		typ, p, err := conn.ReadMessage()
		fmt.Println(typ, "type")
		fmt.Println(string(p), "message")
		fmt.Println(err)
	}
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

func (c *BusinessHandler) UpdatePicture(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {
			name, err := serveUtils.MultipleFileUploadService(r)
			if err != nil {
				erros.JSONError(rw, err, http.StatusBadRequest)
				return
			}

			path := config.Configuration.File

			for _, v := range name {
				_ = os.Rename("./static/"+v, path+bID.String()+"/"+v)
			}

			var newName []string
			for _, v := range name {
				newName = append(newName, config.Configuration.Contents+bID.String()+"/"+v)
			}

			businesses := json_serializer.UpdateAbleBusiness{}
			businesses.Picture = &newName
			err = business2.UpdateBusiness(r, &businesses, bID, "nothing")
			if err != nil {
				erros.JSONError(rw, err, http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return

		} else {
			erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
			return
		}

	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) PatchPicture(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		pid := vars["pid"]

		path := "./static/production/" + bID.String() + "/"
		if config.Configuration.Local == true {
			path = "./static/development/" + bID.String() + "/"
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		present := false
		var ext string
		for _, f := range files {
			name := strings.Split(f.Name(), ".")
			if pid == name[0] {
				present = true
				ext = name[1]
				break
			}
		}

		if present == true {
			if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {
				name, err := serveUtils.FileUploadService(r)
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}
				nameLast := strings.Split(name, ".")

				_ = os.Remove(path + pid + "." + ext)
				_ = os.Rename("./static/"+name, path+pid+"."+nameLast[1])

				var newPics []string
				for _, v := range business.Picture {
					if v != config.Configuration.Contents+bID.String()+"/"+pid+"."+ext {
						newPics = append(newPics, v)
					} else {
						newPics = append(newPics, config.Configuration.Contents+bID.String()+"/"+pid+"."+nameLast[1])
					}
				}

				businesses := json_serializer.UpdateAbleBusiness{}
				businesses.Picture = &newPics
				err = business2.UpdateBusiness(r, &businesses, bID, "nothing")
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}

				rw.WriteHeader(http.StatusNoContent)
				return

			} else {
				erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}
		} else {
			erros.JSONError(rw, errors.New("no resource found"), http.StatusNotFound)
			return
		}
	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) AddPicture(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {
			name, err := serveUtils.MultipleFileUploadService(r)
			if err != nil {
			}

			path := "./static/production/"
			if config.Configuration.Local == true {
				path = "./static/development/"
			}

			for _, v := range name {
				_ = os.Rename("./static/"+v, path+bID.String()+"/"+v)
			}

			for _, v := range name {
				business.Picture = append(business.Picture, config.Configuration.Contents+bID.String()+"/"+v)
			}

			businesses := json_serializer.UpdateAbleBusiness{}
			businesses.Picture = &business.Picture
			err = business2.UpdateBusiness(r, &businesses, bID, "nothing")
			if err != nil {
				erros.JSONError(rw, err, http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return

		} else {
			erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
			return
		}

	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

func (c *BusinessHandler) DeletePicture(rw http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp)
	if userInfo.Authenticated == true {
		vars := mux.Vars(r)
		bID, err := uuid.Parse(vars["id"])
		if err != nil {
			erros.JSONError(rw, err, http.StatusBadRequest)
			return
		}

		service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)
		business, err := service.FindById(bID)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		pid := vars["pid"]

		path := "./static/production/" + bID.String() + "/"
		if config.Configuration.Local == true {
			path = "./static/development/" + bID.String() + "/"
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			erros.JSONError(rw, err, http.StatusNotFound)
			return
		}

		present := false
		var ext string
		for _, f := range files {
			name := strings.Split(f.Name(), ".")
			if pid == name[0] {
				present = true
				ext = name[1]
				break
			}
		}

		if present == true {
			if userInfo.UserId == business.UserId || userInfo.UserType == "admin" || userInfo.UserType == "super" {
				err = os.Remove(path + pid + "." + ext)
				if err != nil {
					erros.JSONError(rw, err, http.StatusNotFound)
					return
				}

				var newPics []string
				for _, v := range business.Picture {
					if v != config.Configuration.Contents+bID.String()+"/"+pid+"."+ext {
						newPics = append(newPics, v)
					}
				}

				businesses := json_serializer.UpdateAbleBusiness{}
				businesses.Picture = &newPics
				err = business2.UpdateBusiness(r, &businesses, bID, "nothing")
				if err != nil {
					erros.JSONError(rw, err, http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusNoContent)
				return

			} else {
				erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}
		} else {
			erros.JSONError(rw, errors.New("no resource found"), http.StatusNotFound)
			return
		}
	}
	erros.JSONError(rw, errors.New("not authorized"), http.StatusUnauthorized)
	return
}

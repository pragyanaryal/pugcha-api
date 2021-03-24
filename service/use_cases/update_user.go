package use_cases

import (
	"errors"
	"fmt"
	"github.com/danhper/structomap"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
	"os"
	"time"
)

func UpdateUser(profile *json_serializer.UpdateAbleUser, id uuid.UUID, content string, r *http.Request) (*json_serializer.UserResponse, error) {
	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
		service     = serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)
	)

	if content == " multipart/form-data" {
		name, err := serveUtils.FileUploadService(r)
		if err != nil {
			return nil, err
		}
		newName := config.Configuration.Contents + name
		profile.ProfilePics = &newName

		go func() {
			_ = os.Rename("./static/"+name, config.Configuration.File + name)
		}()
	}

	err := validateUserRequest(profile)
	if err != nil {
		return nil, err
	}

	profileMap := structomap.New().UseSnakeCase().PickAll().
		OmitIf(func(u interface{}) bool {
			if profile.Status == nil {
				return true
			}
			return false
		}, "Status").
		OmitIf(func(u interface{}) bool {
			if profile.Type == nil {
				return true
			}
			return false
		}, "Type").
		OmitIf(func(u interface{}) bool {
			if profile.FirstName == nil {
				return true
			}
			return false
		}, "FirstName").
		OmitIf(func(u interface{}) bool {
			if profile.LastName == nil {
				return true
			}
			return false
		}, "LastName").
		OmitIf(func(u interface{}) bool {
			if profile.MiddleName == nil {
				return true
			}
			return false
		}, "MiddleName").
		OmitIf(func(u interface{}) bool {
			if profile.Dob == nil {
				return true
			}
			return false
		}, "Dob").
		OmitIf(func(u interface{}) bool {
			if profile.Gender == nil {
				return true
			}
			return false
		}, "Gender").
		OmitIf(func(u interface{}) bool {
			if profile.Contact == nil {
				return true
			}
			return false
		}, "Contact").
		OmitIf(func(u interface{}) bool {
			if profile.ProfilePics == nil {
				return true
			}
			return false
		}, "ProfilePics").
		Transform(profile)

	if len(profileMap) > 0 {
		profileMap["updated_on"] = time.Now()
		err = service.PatchUser(id, &profileMap)
		if err != nil {
			return nil, err
		}
		userById, err := service.GetUserById(id)

		return userById, err
	}

	fmt.Println(len(profileMap), "length update")
	return nil, errors.New("nothing to update")
}

func validateUserRequest(a *json_serializer.UpdateAbleUser) error {
	validate := validator.New()
	_ = validate.RegisterValidation("status", validateStatus)
	_ = validate.RegisterValidation("type", validateType)
	_ = validate.RegisterValidation("gender", validateGender)
	return validate.Struct(a)
}

func validateStatus(fl validator.FieldLevel) bool {
	temp := fl.Field().String()
	if temp != "active" && temp != "blocked" && temp != "approval_needed" && temp != "verification_needed" {
		return false
	}
	return true
}

func validateType(fl validator.FieldLevel) bool {
	temp := fl.Field().String()
	if temp != "user" && temp != "admin" && temp != "super" {
		return false
	}
	return true
}

func validateGender(fl validator.FieldLevel) bool {
	temp := fl.Field().String()
	if temp != "male" && temp != "female" && temp != "rest" && temp != "" {
		return false
	}
	return true
}

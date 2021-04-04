package user

import (
	"errors"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/user_service"
)

func HandleNormalUserOrchestrator(user *json_serializer.CreateNormalUser, rw http.ResponseWriter) (*json_serializer.UserResponse, error) {
	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
		service     = user_service.UserService(userRepo, profileRepo, googleRepo, fbRepo)
	)

	old, err := service.GetUserByEmail(user.Email)

	if err != nil {
		createUser, err := service.CreateUser(user.Email, "user")
		if err != nil {
			return nil, err
		}

		newProfile := createNewProfile(user, createUser.ID)
		err = service.CreateProfile(newProfile)
		if err != nil {
			_ = service.DeleteUser(createUser.ID)
			return nil, err
		}

		if user.SocialPlatform == "google" {
			_, err = service.CreateGoogleProfile(user, createUser.ID, nil)
			if err != nil {
				_ = service.DeleteProfile(createUser.ID)
				_ = service.DeleteUser(createUser.ID)
				return nil, err
			}
		} else {

		}
		newer, err := service.GetUserByEmail(createUser.Email)
		return newer, nil
	}

	if old.Password == "" {
		return old, nil
	}

	return old, errors.New("user already exists")
}

func createNewProfile(user *json_serializer.CreateNormalUser, id uuid.UUID) *models.UserProfile {
	gender := ""
	profile := models.UserProfile{
		UserId:      id,
		Email:       user.Email,
		FirstName:   user.FirstName,
		MiddleName:  nil,
		Gender:      &gender,
		LastName:    user.LastName,
		ProfilePics: user.Picture,
		CreatedOn:   time.Now(),
		UpdatedOn:   time.Now(),
	}
	return &profile
}
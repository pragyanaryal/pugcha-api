package use_cases

import (
	"errors"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"net/http"
	"time"
)

func HandleNormalUserOrchestrator(user *json_serializer.CreateNormalUser, rw http.ResponseWriter) (*json_serializer.UserResponse, error) {
	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
		service     = serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)
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

//func HandleSocialUserOrchestrator(
//	auth string,
//	data *json_serializer.GoogleResponse,
//	fb *json_serializer.FacebookResponse,
//	token *oauth2.Token) (*models.UserProfile, error) {
//
//	var (
//		userRepo    = repositories.UserRepo
//		profileRepo = repositories.UserProfileRepo
//		googleRepo  = repositories.GoogleRepo
//		fbRepo      = repositories.FacebookRepo
//	)
//
//	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)
//
//	if auth == "google" {
//		//	Create User ...
//		user, err := service.CreateUser(data.Email, "user")
//		if err != nil {
//			return nil, err
//		}
//
//		//	Create Profile ...
//		profile := createProfile(user.ID, data, nil)
//		err = service.CreateProfile(profile)
//		if err != nil {
//			//err := service.DeleteUser(user)
//			return nil, err
//		}
//
//		// Create Google Profile ...
//		_, err = service.CreateGoogleProfile(nil, user.ID, token)
//		if err != nil {
//			//err := service.DeleteProfile(profile)
//			//err = service.DeleteUser(user)
//			return nil, err
//		}
//
//		return profile, nil
//
//	} else {
//
//		return nil, nil
//	}
//}
//
//func createProfile(
//	id uuid.UUID,
//	dat *json_serializer.GoogleResponse,
//	fb *json_serializer.FacebookResponse) *models.UserProfile {
//
//	if dat != nil {
//		gender := ""
//		profile := models.UserProfile{
//			UserId:      id,
//			Email:       dat.Email,
//			FirstName:   dat.Name,
//			MiddleName:  nil,
//			Gender:      &gender,
//			LastName:    dat.FamilyName,
//			ProfilePics: dat.Picture,
//			CreatedOn:   time.Now(),
//			UpdatedOn:   time.Now(),
//		}
//		return &profile
//	}
//	return nil
//}

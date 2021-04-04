package user

import (
	"net/http"
	
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/user_service"
)

func GoogleCallBackOrchestrator(r *http.Request) (*models.UserProfile, error) {
	token, data, err := user_service.GoogleService.GoogleCallback(r)
	if err != nil {
		return nil, err
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := user_service.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	_, err = service.GetUserByEmail(data.Email)

	//if err != nil {
	//	prof, err := HandleSocialUserOrchestrator("google", data, nil, token)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return prof, nil
	//}
	profile, err := UpdateSocialUserOrchestrator(data, nil, token)
	return profile, nil
}

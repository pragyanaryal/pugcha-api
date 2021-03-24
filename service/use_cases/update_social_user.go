package use_cases

import (
	"errors"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"golang.org/x/oauth2"
	"time"
)

func UpdateSocialUserOrchestrator(
	data *json_serializer.GoogleResponse,
	fb *json_serializer.FacebookResponse,
	token *oauth2.Token) (*models.UserProfile, error) {

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	if data != nil {
		profile, err := service.GetProfileByEmail(data.Email)
		if err != nil {
			return nil, err
		}

		profile, change := updateLatestProfile(profile, data, nil)
		if change == true {
			err = service.UpdateProfile(profile)
		}

		googleProfile, err := service.GetGoogleProfileByEmail(data.Email)
		if err != nil {
			return profile, nil
		}

		googleProfile = updateGoogleProfile(googleProfile, token)
		err = service.UpdateGoogleProfile(googleProfile)
		if err != nil {
		}
		return profile, nil
	}
	return nil, errors.New("not found")
}

func updateGoogleProfile(account *models.GoogleAccount, token *oauth2.Token) *models.GoogleAccount {
	account.AccessKey = token.AccessToken
	account.RefreshKey = token.RefreshToken

	return account
}

func updateLatestProfile(
	profile *models.UserProfile,
	data *json_serializer.GoogleResponse,
	fb *json_serializer.FacebookResponse) (*models.UserProfile, bool) {

	if data != nil {
		comp1 := profile.ProfilePics == data.Picture
		profile.ProfilePics = data.Picture

		comp2 := profile.LastName == data.FamilyName
		profile.LastName = data.FamilyName

		comp3 := profile.FirstName == data.GivenName
		profile.FirstName = data.GivenName

		if comp1 == comp2 == comp3 == true {
			return profile, false
		}

		profile.UpdatedOn = time.Now()

		return profile, true
	}

	return nil, false
}

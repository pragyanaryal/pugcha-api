package user_service

import (
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"golang.org/x/oauth2"
)

// Create Google profile of the user
func (users *userService) CreateGoogleProfile(data *json_serializer.CreateNormalUser, id uuid.UUID, token *oauth2.Token) (
	*models.GoogleAccount, error) {

	googleProfile := models.GoogleAccount{
		UserId:   id,
		Email:    data.Email,
		GoogleID: data.SocialID,
	}
	if token != nil {
		googleProfile.AccessKey = token.AccessToken
		googleProfile.RefreshKey = token.RefreshToken
	}

	profileMap := structomap.New().UseSnakeCase().PickAll().Transform(googleProfile)

	err := users.googleProfile.CreateGoogleProfile(&profileMap)
	if err != nil {
		return nil, err
	}
	return &googleProfile, nil
}

// Check if google profile exists
func (users *userService) GetGoogleProfileByEmail(email string) (*models.GoogleAccount, error) {
	profile, err := users.googleProfile.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Update google profile of the user in the system
func (users *userService) UpdateGoogleProfile(account *models.GoogleAccount) error {
	err := users.googleProfile.UpdateGoogleProfile(account)
	if err != nil {
		return err
	}

	return nil
}

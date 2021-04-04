package user_service

import (
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
)

// Create Profile for the user ...
func (users *userService) CreateProfile(profile *models.UserProfile) error {
	profileMap := structomap.New().UseSnakeCase().PickAll().
		OmitIf(func(u interface{}) bool {
			if profile.Dob == nil {
				return true
			}
			return false
		}, "Dob").
		OmitIf(func(u interface{}) bool {
			if profile.Contact == nil {
				return true
			}
			return false
		}, "Contact").
		OmitIf(func(u interface{}) bool {
			if profile.Gender == nil {
				return true
			}
			return false
		}, "Gender").
		Transform(profile)

	err := users.userProfile.CreateProfile(&profileMap)
	if err != nil {
		return err
	}

	return err
}

// Get profile of the user by email
func (users *userService) GetProfiles() ([]*models.UserProfile, error) {
	profile, err := users.userProfile.ListProfile()
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// Get profile of the user by email
func (users *userService) GetProfileByEmail(email string) (*models.UserProfile, error) {
	profile, err := users.userProfile.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// Get profile of the user by id
func (users *userService) GetProfileById(id uuid.UUID) (*models.UserProfile, error) {
	profile, err := users.userProfile.FindById(id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// Update Profile for the user ...
func (users *userService) UpdateProfile(profile *models.UserProfile) error {
	err := users.userProfile.UpdateProfile(profile)
	if err != nil {
		return err
	}
	return nil
}

// Update Profile for the user ...
func (users *userService) PatchProfile(profile *models.UserProfile, prof *map[string]interface{}) error {
	err := users.userProfile.PatchProfile(profile.UserId, prof)
	if err != nil {
		return err
	}
	return nil
}

// Delete Profile of the user ...
func (users *userService) DeleteProfile(profile uuid.UUID) error {
	err := users.userProfile.DeleteProfile(profile)
	if err != nil {
		return err
	}
	return nil
}

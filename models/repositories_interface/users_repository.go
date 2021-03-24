package repositories_interface

import (
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
)

type UserRepository interface {
	FindByEmail(string) (*json_serializer.UserResponse, error)
	FindById(uuid.UUID) (*json_serializer.UserResponse, error)
	CreateUser(*map[string]interface{}) error
	PatchUser(uuid.UUID, *map[string]interface{}) error
	ListUser(filter *serveUtils.SQLFilter) (*[]*json_serializer.UserResponse, error)
	DeleteUser(uuid.UUID) error
	QueryUser(a ...interface{}) ([]*models.User, error)
	UpdateUser(*models.User) error
}

type UserProfileRepository interface {
	FindByEmail(string) (*models.UserProfile, error)
	FindById(uuid.UUID) (*models.UserProfile, error)
	CreateProfile(*map[string]interface{}) error
	PatchProfile(uuid.UUID, *map[string]interface{}) error
	ListProfile() ([]*models.UserProfile, error)
	DeleteProfile(uuid.UUID) error
	QueryProfile(a ...interface{}) ([]*models.UserProfile, error)
	UpdateProfile(user *models.UserProfile) error
}

type GoogleProfile interface {
	FindByEmail(email string) (*models.GoogleAccount, error)
	FindById(id uuid.UUID) (*models.GoogleAccount, error)
	ListGoogleProfile() ([]*models.GoogleAccount, error)
	CreateGoogleProfile(user *map[string]interface{}) error
	PatchProfile(prof *models.GoogleAccount, user *map[string]interface{}) error
	UpdateGoogleProfile(user *models.GoogleAccount) error
	DeleteGoogleProfile(user *models.GoogleAccount) error
}

type FacebookProfile interface {
	FindByEmail(email string) (*models.FacebookAccount, error)
	FindById(id uuid.UUID) (*models.FacebookAccount, error)
	ListFacebookProfile() ([]*models.FacebookAccount, error)
	CreateFacebookProfile(user *models.FacebookAccount) error
	UpdateFacebookProfile(user *models.FacebookAccount) error
	DeleteFacebookProfile(user *models.FacebookAccount) error
}

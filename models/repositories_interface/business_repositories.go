package repositories_interface

import (
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
)

type BusinessRepository interface {
	FindById(uuid.UUID) (*json_serializer.SmallBusinessResponse, error)
	ListBusiness(filter *serveUtils.SQLFilter) (*[]interface{}, error)
	CreateBusiness(*map[string]interface{}) error
	DeleteBusiness(uuid.UUID) error
	PatchBusiness(uuid.UUID, *map[string]interface{}) error
	DeleteAddress(uuid.UUID) error
	BatchDeleteAddress([]uuid.UUID) error
}

type BusinessProfileRepository interface {
	CreateProfile(*map[string]interface{}) error
	FindById(uuid.UUID) (*json_serializer.FullBusinessResponse, error)
	ListBusinessProfile() ([]*models.BusinessProfile, error)
	PatchProfile(uuid.UUID, *map[string]interface{}) error
	DeleteProfile(uuid.UUID) error
}

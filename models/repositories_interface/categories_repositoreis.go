package repositories_interface

import (
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
)

type CategoriesRepository interface {
	FindById(uuid.UUID) (*models.Categories, error)
	ListCategories() ([]*models.Categories, error)
	CreateCategories(*map[string]interface{}) error
	PatchCategories(uuid.UUID, *map[string]interface{}) error
	DeleteCategories(uuid.UUID) error
}

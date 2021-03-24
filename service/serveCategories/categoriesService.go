package serveCategories

import (
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"time"
)

// Create Category
func (category *categoryService) CreateCategory(name string, picture string, id uuid.UUID) (*models.Categories, error) {
	ids := uuid.New()
	categories := models.Categories{
		Id:        ids,
		Name:      name,
		CreatedBy: id,
		Picture:   config.Configuration.Contents + ids.String() + "/" + picture,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}

	categoryMap := structomap.New().UseSnakeCase().PickAll().Transform(categories)

	err := category.categoryRepo.CreateCategories(&categoryMap)
	if err != nil {
		return nil, err
	}
	return &categories, nil
}

// List Categories
func (category *categoryService) GetCategories() ([]*models.Categories, error) {
	categories, err := category.categoryRepo.ListCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Get Category
func (category *categoryService) GetCategory(id uuid.UUID) (*models.Categories, error) {
	categories, err := category.categoryRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update Categories
func (category *categoryService) UpdateCategories(cat *models.Categories, prof *map[string]interface{}) error {
	err := category.categoryRepo.PatchCategories(cat.Id, prof)
	if err != nil {
		return err
	}
	return nil
}

// Delete Categories
func (category *categoryService) DeleteCategory(categories *models.Categories) error {
	err := category.categoryRepo.DeleteCategories(categories.Id)
	if err != nil {
		return err
	}
	return nil
}

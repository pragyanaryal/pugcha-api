package serveCategories

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/models/repositories_interface"
)

type categoryService struct {
	categoryRepo repositories_interface.CategoriesRepository
}

func CategoriesService(categoryRepo repositories_interface.CategoriesRepository) *categoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

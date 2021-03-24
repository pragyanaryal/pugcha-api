package use_cases

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveCategories"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
	"os"
)

func CreateCategoriesOrchestrator(r *http.Request, createdBy uuid.UUID) (*models.Categories, error) {
	name, err := serveUtils.FileUploadService(r)
	if err != nil {
		return nil, err
	}

	service := serveCategories.CategoriesService(repositories.CategoryRepo)

	if val, ok := r.Form["name"]; ok {
		categories := json_serializer.CreateCategoriesRequest{
			Name:      val[0],
			Picture:   name,
			CreatedBy: createdBy,
		}

		err = validateRequest(&categories)
		if err != nil {
			return nil, err
		}

		category, err := service.CreateCategory(categories.Name, categories.Picture, categories.CreatedBy)
		if err != nil {
			return nil, err
		}

		go func() {
			_, err = os.Stat(category.Id.String())
			if os.IsNotExist(err) {
				_ = os.MkdirAll(config.Configuration.File + category.Id.String(), 0755)
					_ = os.Rename("./static/"+name, config.Configuration.File + category.Id.String() + "/"+ name)
			}
		}()


		return category, nil
	}

	return nil, errors.New("name not included")
}

func validateRequest(request *json_serializer.CreateCategoriesRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}

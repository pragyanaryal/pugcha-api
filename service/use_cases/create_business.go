package use_cases

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveBusiness"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
	"os"
)

func CreateBusinessOrchestrator(r *http.Request) (*models.BusinessProfile, error) {
	userId := r.Context().Value(middleware.AuthUserTemp{}).(middleware.AuthUserTemp).UserId

	name, err := serveUtils.MultipleFileUploadService(r)
	if err != nil {
		return nil, err
	}

	profile, err := serveBusiness.GetBusinessProfileData(r, userId, name)
	if err != nil {
		return nil, err
	}

	err = validateBusinessRequest(profile)
	if err != nil {
		return nil, err
	}

	service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	// create business
	business, err := service.CreateBusiness(userId)
	if err != nil {
		return nil, err
	}

	var newName []string
	for _, v := range name {
		newName = append(newName, config.Configuration.Contents+business.ID.String()+"/"+v)
	}
	profile.Picture = newName

	go func() {
		_, err = os.Stat(business.ID.String())
		if os.IsNotExist(err) {
			_ = os.MkdirAll(config.Configuration.File + business.ID.String(), 0755)
			for _, v := range name {
				_ = os.Rename("./static/"+v, config.Configuration.File+business.ID.String()+"/"+v)
			}
		}
	}()

	profile, _ = serveBusiness.PopulateBusinessId(business.ID, profile)

	// create profile for business
	err = service.CreateBusinessProfile(profile)

	if err != nil {
		_ = service.DeleteBusiness(business)
		for _, v := range name {
			_ = os.Remove("./static/" + v)
		}
		return nil, err
	}
	return profile, nil
}

func validateBusinessRequest(a *models.BusinessProfile) error {
	validate := validator.New()
	return validate.Struct(a)
}

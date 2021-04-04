package business

import (
	"net/http"
	"time"
	
	"github.com/danhper/structomap"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/business_service"
)

func UpdateBusiness(r *http.Request, profile *json_serializer.UpdateAbleBusiness, id uuid.UUID, content string) error {

	err := validateBusiness(profile)
	if err != nil {
		return err
	}

	profile, err = business_service.PopulateBusinessIdOnUpdate(id, profile)
	if err != nil {
		return err
	}

	profileMap := structomap.New().UseSnakeCase().PickAll().
		OmitIf(func(u interface{}) bool {
			if profile.Approved == nil {
				return true
			}
			return false
		}, "Approved").
		OmitIf(func(u interface{}) bool {
			if profile.UserId == nil {
				return true
			}
			return false
		}, "UserId").
		OmitIf(func(u interface{}) bool {
			if profile.ApprovedBy == nil {
				return true
			}
			return false
		}, "ApprovedBy").
		OmitIf(func(u interface{}) bool {
			if profile.Blocked == nil {
				return true
			}
			return false
		}, "Blocked").
		OmitIf(func(u interface{}) bool {
			if profile.CategoryId == nil {
				return true
			}
			return false
		}, "CategoryId").
		OmitIf(func(u interface{}) bool {
			if profile.Name == nil {
				return true
			}
			return false
		}, "Name").
		OmitIf(func(u interface{}) bool {
			if profile.PanNumber == nil {
				return true
			}
			return false
		}, "PanNumber").
		OmitIf(func(u interface{}) bool {
			if profile.Email == nil {
				return true
			}
			return false
		}, "Email").
		OmitIf(func(u interface{}) bool {
			if profile.Website == nil {
				return true
			}
			return false
		}, "Website").
		OmitIf(func(u interface{}) bool {
			if profile.Picture == nil {
				return true
			}
			return false
		}, "Picture").
		OmitIf(func(u interface{}) bool {
			if profile.VatNumber == nil {
				return true
			}
			return false
		}, "VatNumber").
		OmitIf(func(u interface{}) bool {
			if profile.Contact == nil {
				return true
			}
			return false
		}, "Contact").
		OmitIf(func(u interface{}) bool {
			if profile.Location == nil {
				return true
			}
			return false
		}, "Location").
		OmitIf(func(u interface{}) bool {
			if profile.EstablishedDate == nil {
				return true
			}
			return false
		}, "EstablishedDate").
		OmitIf(func(u interface{}) bool {
			if profile.Owner == nil {
				return true
			}
			return false
		}, "Owner").
		OmitIf(func(u interface{}) bool {
			if profile.Address == nil {
				return true
			}
			return false
		}, "Address").
		OmitIf(func(u interface{}) bool {
			if profile.OpeningHours == nil {
				return true
			}
			return false
		}, "OpeningHours").
		OmitIf(func(u interface{}) bool {
			if profile.Description == nil {
				return true
			}
			return false
		}, "Description").
		Transform(profile)

	service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	if len(profileMap) > 0 {
		profileMap["updated_on"] = time.Now()

		err = service.PatchBusiness(id, &profileMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateBusiness(business *json_serializer.UpdateAbleBusiness) error {
	validate := validator.New()
	return validate.Struct(business)
}

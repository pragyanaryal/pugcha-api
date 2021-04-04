package business_service

import (
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"time"
)

func (businesses *businessService) FindProfileById(id uuid.UUID) (*json_serializer.FullBusinessResponse, error) {
	business, err := businesses.profileRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (businesses *businessService) ListBusinessProfiles() ([]*models.BusinessProfile, error) {
	bus, err := businesses.profileRepo.ListBusinessProfile()
	if err != nil {
		return nil, err
	}
	return bus, nil
}

func (businesses *businessService) CreateBusinessProfile(busies *models.BusinessProfile) error {

	busies.CreatedOn = time.Now().UTC()
	busies.UpdatedOn = time.Now().UTC()

	businessMap := structomap.New().UseSnakeCase().PickAll().
		Omit("Geom").
		OmitIf(func(u interface{}) bool {
			if busies.EstablishedDate == nil {
				return true
			}
			return false
		}, "EstablishedDate").
		OmitIf(func(u interface{}) bool {
			if busies.PanNumber == nil {
				return true
			}
			return false
		}, "PanNumber").
		OmitIf(func(u interface{}) bool {
			if busies.Email == nil {
				return true
			}
			return false
		}, "Email").
		OmitIf(func(u interface{}) bool {
			if busies.Website == nil {
				return true
			}
			return false
		}, "Website").
		OmitIf(func(u interface{}) bool {
			if busies.VatNumber == nil {
				return true
			}
			return false
		}, "VatNumber").
		OmitIf(func(u interface{}) bool {
			if busies.Description == nil {
				return true
			}
			return false
		}, "Description").
		OmitIf(func(u interface{}) bool {
			if busies.Location.Latitude == nil || busies.Location.Longitude == nil {
				return true
			}
			return false
		}, "Location").
		Transform(busies)

	err := businesses.profileRepo.CreateProfile(&businessMap)
	if err != nil {
		return err
	}
	return nil
}

func (businesses *businessService) PatchBusinessProfile(bu *models.BusinessProfile, bus *map[string]interface{}) error {
	err := businesses.profileRepo.PatchProfile(bu.BusinessId, bus)
	if err != nil {
		return err
	}
	return nil
}

func (businesses *businessService) DeleteBusinessProfile(business *models.BusinessProfile) error {
	err := businesses.profileRepo.DeleteProfile(business.BusinessId)
	if err != nil {
		return err
	}
	return nil
}

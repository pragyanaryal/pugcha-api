package serveBusiness

import (
	"fmt"
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"time"
)

func (businesses *businessService) FindById(id uuid.UUID) (*json_serializer.SmallBusinessResponse, error) {
	business, err := businesses.businessRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (businesses *businessService) ListBusinesses(filter *serveUtils.SQLFilter) (*[]interface{}, error) {
	bus, err := businesses.businessRepo.ListBusiness(filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bus, nil
}

func (businesses *businessService) CreateBusiness(userId uuid.UUID) (*models.Business, error) {
	company := models.Business{
		ID:        uuid.New(),
		Approved:  false,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Blocked:   false,
		UserId:    userId,
	}
	businessMap := structomap.New().UseSnakeCase().PickAll().
		OmitIf(func(u interface{}) bool {
			if company.ApprovedBy == uuid.Nil {
				return true
			}
			return false
		}, "ApprovedBy").
		Transform(company)

	err := businesses.businessRepo.CreateBusiness(&businessMap)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (businesses *businessService) DeleteBusiness(business *models.Business) error {
	err := businesses.businessRepo.DeleteBusiness(business.ID)
	if err != nil {
		return err
	}
	return nil
}

func (businesses *businessService) DeleteAddressById(id uuid.UUID) error {
	err := businesses.businessRepo.DeleteAddress(id)
	if err != nil {
		return err
	}
	return nil
}

func (businesses *businessService) BatchDeleteAddressById(id []uuid.UUID) error {
	err := businesses.businessRepo.BatchDeleteAddress(id)
	if err != nil {
		return err
	}
	return nil
}

func (businesses *businessService) PatchBusiness(business uuid.UUID, bus *map[string]interface{}) error {
	err := businesses.businessRepo.PatchBusiness(business, bus)
	if err != nil {
		return err
	}
	return nil
}

package repositories

import (
	"errors"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"

	"github.com/google/uuid"
)

var (
	FacebookRepo = &facebookRepo{}
)

type facebookRepo struct{}

func (user *facebookRepo) FindByEmail(email string) (*models.FacebookAccount, error) {
	return nil, errors.New("not found")
}

func (user *facebookRepo) FindById(id uuid.UUID) (*models.FacebookAccount, error) {
	return nil, errors.New("not found")
}

func (user *facebookRepo) ListFacebookProfile() ([]*models.FacebookAccount, error) {
	return nil, nil
}

func (user *facebookRepo) CreateFacebookProfile(profile *models.FacebookAccount) error {
	return nil
}

func (user *facebookRepo) UpdateFacebookProfile(profile *models.FacebookAccount) error {
	return nil
}

func (user *facebookRepo) DeleteFacebookProfile(profile *models.FacebookAccount) error {
	return nil
}

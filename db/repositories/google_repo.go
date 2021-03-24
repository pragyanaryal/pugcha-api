package repositories

import (
	"context"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"

	"github.com/google/uuid"
)

var GoogleRepo = &googleRepo{}

type googleRepo struct{}

func (user *googleRepo) FindByEmail(email string) (*models.GoogleAccount, error) {
	sql, args, err := sqrl.Select("*").From("google_profile").Where(sqrl.Eq{"email": email}).
		PlaceholderFormat(sqrl.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var user2 models.GoogleAccount
	err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return &user2, nil
}

func (user *googleRepo) FindById(id uuid.UUID) (*models.GoogleAccount, error) {
	sql, args, err := sqrl.Select("*").From("google_profile").Where(sqrl.Eq{"user_id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var user2 models.GoogleAccount
	err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return &user2, nil
}

func (user *googleRepo) CreateGoogleProfile(profile *map[string]interface{}) error {
	sql, args, err := sqrl.Insert("google_profile").
		SetMap(*profile).PlaceholderFormat(sqrl.Dollar).ToSql()

	if err != nil {
		return err
	}

	exec, err := db.Pool.Exec(context.Background(), sql, args[:]...)
	if err != nil {
		return err
	}

	count := exec.RowsAffected()
	if count != 1 {
		return err
	}

	return nil
}

func (user *googleRepo) PatchProfile(prof *models.GoogleAccount, profile *map[string]interface{}) error {
	sql, args, err := sqrl.Update("google_profile").
		SetMap(*profile).
		Where(sqrl.Eq{"user_id": prof.UserId}).
		PlaceholderFormat(sqrl.Dollar).ToSql()

	if err != nil {
		return err
	}

	exec, err := db.Pool.Exec(context.Background(), sql, args[:]...)
	if err != nil {
		return err
	}

	count := exec.RowsAffected()
	if count != 1 {
		return err
	}

	return nil
}

func (user *googleRepo) DeleteGoogleProfile(profile *models.GoogleAccount) error {
	sql, args, err := sqrl.Delete("google_profile").
		Where(sqrl.Eq{"user_id": profile.UserId}).PlaceholderFormat(sqrl.Dollar).ToSql()

	if err != nil {
		return err
	}

	exec, err := db.Pool.Exec(context.Background(), sql, args[:]...)
	if err != nil {
		return err
	}

	count := exec.RowsAffected()
	if count != 1 {
		return err
	}

	return nil
}

func (user *googleRepo) ListGoogleProfile() ([]*models.GoogleAccount, error) {
	sql, args, err := sqrl.Select("*").From("google_profile").PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 []*models.GoogleAccount
	err = pgxscan.Select(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return user2, nil
}

func (user *googleRepo) UpdateGoogleProfile(profile *models.GoogleAccount) error {
	return nil
}

package repositories

import (
	"context"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
)

var UserProfileRepo = &userProfileRepo{}

type userProfileRepo struct{}

func (user *userProfileRepo) FindByEmail(email string) (*models.UserProfile, error) {
	sql, args, err := sqrl.Select("*").From("users_profile").Where(sqrl.Eq{"email": email}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 models.UserProfile
	err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return &user2, nil
}

func (user *userProfileRepo) FindById(id uuid.UUID) (*models.UserProfile, error) {
	sql, args, err := sqrl.Select("*").From("users_profile").Where(sqrl.Eq{"user_id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 models.UserProfile
	err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return &user2, nil
}

func (user *userProfileRepo) CreateProfile(profile *map[string]interface{}) error {
	sql, args, err := sqrl.Insert("users_profile").
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

func (user *userProfileRepo) PatchProfile(userId uuid.UUID, profile *map[string]interface{}) error {
	sql, args, err := sqrl.Update("users_profile").
		SetMap(*profile).
		Where(sqrl.Eq{"user_id": userId}).
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

func (user *userProfileRepo) ListProfile() ([]*models.UserProfile, error) {
	sql, args, err := sqrl.Select("*").From("users_profile").PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 []*models.UserProfile
	err = pgxscan.Select(context.Background(), db.Pool, &user2, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return user2, nil
}

func (user *userProfileRepo) DeleteProfile(id uuid.UUID) error {
	sql, args, err := sqrl.Delete("users_profile").
		Where(sqrl.Eq{"user_id": id}).
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

func (user *userProfileRepo) QueryProfile(a ...interface{}) ([]*models.UserProfile, error) {
	return nil, nil
}

func (user *userProfileRepo) UpdateProfile(*models.UserProfile) error {
	return nil
}

package repositories

import (
	"context"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
)

var CategoryRepo = &categoryRepo{}

type categoryRepo struct{}

func (categories *categoryRepo) FindById(id uuid.UUID) (*models.Categories, error) {
	sql, args, err := sqrl.Select("*").From("categories").Where(sqrl.Eq{"id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var category models.Categories
	if err = pgxscan.Get(context.Background(), db.Pool, &category, sql, args[:]...); err != nil {
		return nil, err
	}

	return &category, nil
}

func (categories *categoryRepo) ListCategories() ([]*models.Categories, error) {
	sql, _, err := sqrl.Select("*").From("categories").PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var category []*models.Categories
	if err = pgxscan.Select(context.Background(), db.Pool, &category, sql); err != nil {
		return nil, err
	}

	return category, nil
}

func (categories *categoryRepo) CreateCategories(category *map[string]interface{}) error {
	sql, args, err := sqrl.Insert("categories").
		SetMap(*category).
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

func (categories *categoryRepo) PatchCategories(id uuid.UUID, prof *map[string]interface{}) error {
	sql, args, err := sqrl.Update("categories").
		SetMap(*prof).
		Where(sqrl.Eq{"id": id}).
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

func (categories *categoryRepo) DeleteCategories(id uuid.UUID) error {
	sql, args, err := sqrl.Delete("categories").
		Where(sqrl.Eq{"id": id}).PlaceholderFormat(sqrl.Dollar).ToSql()

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

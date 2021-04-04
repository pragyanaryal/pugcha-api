package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"strconv"
	"strings"
	"time"
)

var UserRepo = &userRepo{}

type userRepo struct{}

var column = []string{"id", "status", "password", "email", "type", "profile"}

func (user *userRepo) FindByEmail(email string) (*json_serializer.UserResponse, error) {
	sql, args, err := sqrl.Select(strings.Join(column, ",")).From("users_view").Where(sqrl.Eq{"email": email}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 json_serializer.UserResponse
	if err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...); err != nil {
		return nil, errors.New("not found")
	}

	return &user2, nil
}

func (user *userRepo) FindById(id uuid.UUID) (*json_serializer.UserResponse, error) {
	sql, args, err := sqrl.Select(strings.Join(column, ",")).From("users_view").Where(sqrl.Eq{"id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var user2 json_serializer.UserResponse
	if err = pgxscan.Get(context.Background(), db.Pool, &user2, sql, args[:]...); err != nil {
		return nil, errors.New("not found")
	}

	return &user2, nil
}

func (user *userRepo) CreateUser(users *map[string]interface{}) error {
	sql, args, err := sqrl.Insert("users").SetMap(*users).PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return err
	}

	exec, err := db.Pool.Exec(context.Background(), sql, args[:]...)
	if err != nil {
		return err
	}

	if exec.RowsAffected() != 1 {
		return errors.New("not inserted")
	}

	if val, ok := (*users)["user_type"]; ok {
		if val != "user"{
			sql, args, err = sqrl.Insert("users_profile").
				Columns("user_id").Values((*users)["id"]).
				PlaceholderFormat(sqrl.Dollar).ToSql()

			exec, err = db.Pool.Exec(context.Background(), sql, args[:]...)
			if err != nil {
				return err
			}

			if exec.RowsAffected() != 1 {
				return errors.New("not inserted")
			}
		}
	}

	return nil
}

func (user *userRepo) PatchUser(id uuid.UUID, patch *map[string]interface{}) error {
	temp := map[string]interface{}{}
	if val, ok := (*patch)["status"]; ok {
		temp["status"] = val
		delete(*patch, "status")
	}
	if val, ok := (*patch)["type"]; ok {
		temp["type"] = val
		delete(*patch, "type")
	}
	if val, ok := (*patch)["password"]; ok {
		temp["password"] = val
		delete(*patch, "password")
	}
	if _, ok := (*patch)["updated_on"]; ok {
		delete(*patch, "updated_on")
	}

	b1 := pgx.Batch{}
	if len(temp) > 0 {
		temp["updated_on"] = time.Now()
		sql, args, err := sqrl.Update("users").
			SetMap(temp).
			Where(sqrl.Eq{"id": id}).
			PlaceholderFormat(sqrl.Dollar).ToSql()

		if err != nil {
			return err
		}

		b1.Queue(sql, args[:]...)
	}

	if len(*patch) > 0 {
		(*patch)["updated_on"] = time.Now()
		sql, arg, errs := sqrl.Update("users_profile").
			SetMap(*patch).
			Where(sqrl.Eq{"user_id": id}).
			PlaceholderFormat(sqrl.Dollar).ToSql()
		if errs != nil {
			return errs
		}
		b1.Queue(sql, arg[:]...)
	}

	batch := db.Pool.SendBatch(context.Background(), &b1)
	fmt.Println(b1.Len(), "lenss")

	for i := 0; i < b1.Len(); i++ {
		exec, err := batch.Exec()
		if err != nil {
			_ = batch.Close()
			return err
		}
		if exec.RowsAffected() != 1 {
			fmt.Println(exec.RowsAffected(), "effect")
			_ = batch.Close()
			return errors.New("problem inserting")
		}
	}

	_ = batch.Close()
	return nil
}

func (user *userRepo) ListUser(filter *serveUtils.SQLFilter) (*[]*json_serializer.UserResponse, error) {
	sql, _, err := BuildUserListQuery(filter)
	if err != nil {
		return nil, err
	}

	var user2 []*json_serializer.UserResponse

	if err = pgxscan.Select(context.Background(), db.Pool, &user2, sql); err != nil {
		return nil, err
	}
	return &user2, nil
}

func (user *userRepo) DeleteUser(id uuid.UUID) error {
	sql, args, err := sqrl.Delete("users").
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

func (user *userRepo) QueryUser(a ...interface{}) ([]*models.User, error) {
	return nil, nil
}

func (user *userRepo) UpdateUser(*models.User) error {
	return nil
}

func BuildUserListQuery(params *serveUtils.SQLFilter) (sql string, args []interface{}, err error) {
	var tempSort []string

	for _, val := range params.Sort {
		tempSort = append(tempSort, val)
	}
	whereStarted := false

	n := "SELECT " + strings.Join(column, ",") + " FROM users_view "

	if _, ok := params.Filter["status"]; ok {
		n = n + "WHERE status = '" + params.Args["status"][0].(string) + "' "
		delete(params.Filter, "status")
		delete(params.Args, "status")
		whereStarted = true
	}
	if _, ok := params.Filter["type"]; ok {
		if whereStarted == false {
			n = n + "WHERE type = '" + params.Args["type"][0].(string) + "' "
		} else {
			n = n + "AND type = '" + params.Args["type"][0].(string) + "' "
		}
		delete(params.Filter, "type")
		delete(params.Args, "type")
		whereStarted = true
	}
	if _, ok := params.Filter["gender"]; ok {
		if whereStarted == false {
			n = n + "WHERE gender = '" + params.Args["gender"][0].(string) + "' "
		} else {
			n = n + "AND gender = '" + params.Args["gender"][0].(string) + "' "
		}
		delete(params.Filter, "gender")
		delete(params.Args, "gender")
		whereStarted = true
	}

	keys := make([]string, 0, len(params.Filter))

	for k := range params.Filter {
		keys = append(keys, k+" = "+params.Args[k][0].(string))
	}

	if len(keys) > 0 {
		if whereStarted == false {
			n = n + fmt.Sprintf("WHERE %s ", strings.Join(keys, " AND "))
		} else {
			n = n + fmt.Sprintf("AND %s ", strings.Join(keys, " AND "))
		}
	}

	if len(tempSort) > 0 {
		n = n + fmt.Sprintf("ORDER BY %s ", strings.Join(tempSort, ","))
	}

	n = n + "LIMIT " + strconv.Itoa(int(params.Limit)) + " OFFSET " + strconv.Itoa(int(params.Offset))
	return n, nil, nil
}

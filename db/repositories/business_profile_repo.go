package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/danhper/structomap"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/log"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"strconv"
	"strings"
)

var BusinessProfileRepo = &businessProfileRepo{}

type businessProfileRepo struct{}

var tempNewWhere = []string{"picture", "business_id", "contact", "location", "description",
	"category_name", "business_name", "pan_number", "vat_number", "email", "website", "established_date", "address", "opening"}

func (businesses *businessProfileRepo) FindById(id uuid.UUID) (*json_serializer.FullBusinessResponse, error) {
	tempCol := tempNewWhere
	tempCol = append(tempCol, "user_id")

	sql, args, err := sqrl.Select(strings.Join(tempCol, ",")).From("businesses_view").Where(sqrl.Eq{"business_id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	oids := pgx.QueryResultFormatsByOID{db.GeoOid: pgx.BinaryFormatCode}

	temp := args[0]
	args[0] = oids

	args = append(args, temp)

	var business json_serializer.FullBusinessResponse
	if err = pgxscan.Get(context.Background(), db.Pool, &business, sql, args[:]...); err != nil {
		return nil, err
	}

	return &business, nil
}

func (businesses *businessProfileRepo) ListBusinessProfile() ([]*models.BusinessProfile, error) {
	sql, _, err := sqrl.Select("*").
		From("business_profiles").Limit(50).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var business []*models.BusinessProfile

	oids := pgx.QueryResultFormatsByOID{db.GeoOid: pgx.BinaryFormatCode}

	err = pgxscan.Select(context.Background(), db.Pool, &business, sql, oids)
	if err != nil {
		return nil, err
	}

	return business, nil
}

func (businesses *businessProfileRepo) CreateProfile(business *map[string]interface{}) error {
	locate := models.Location{}
	var address []models.Address
	var owner []models.Owner
	var openingHours []models.OpeningHours

	_ = mapstructure.Decode((*business)["location"], &locate)
	_ = mapstructure.Decode((*business)["address"], &address)
	_ = mapstructure.Decode((*business)["owner"], &owner)
	_ = mapstructure.Decode((*business)["opening_hours"], &openingHours)

	delete(*business, "location")
	delete(*business, "address")
	delete(*business, "owner")
	delete(*business, "opening_hours")

	var sql string
	var args []interface{}

	if locate.Latitude != nil && locate.Longitude != nil {
		sql, args = businesses.BuildQuery(business, &locate)
	} else {
		sql, args = businesses.BuildQuery(business, nil)
	}

	b1 := pgx.Batch{}
	b1.Queue(sql, args[:]...)

	for _, val := range address {
		addMap := structomap.New().UseSnakeCase().PickAll().Transform(val)
		sqls, argss, err := sqrl.Insert("addresses").SetMap(addMap).PlaceholderFormat(sqrl.Dollar).ToSql()
		if err != nil {
			log.Err(err)
			continue
		}
		b1.Queue(sqls, argss[:]...)
	}
	for _, val := range owner {
		addMap := structomap.New().UseSnakeCase().PickAll().Transform(val)
		sql, args, err := sqrl.Insert("owners_info").SetMap(addMap).PlaceholderFormat(sqrl.Dollar).ToSql()
		if err != nil {
			log.Err(err)
			continue
		}
		b1.Queue(sql, args[:]...)
	}
	for _, val := range openingHours {
		addMap := structomap.New().UseSnakeCase().PickAll().Transform(val)

		sqls, argss, err := sqrl.Insert("opening_times").SetMap(addMap).PlaceholderFormat(sqrl.Dollar).ToSql()
		if err != nil {
			log.Err(err)
			continue
		}
		b1.Queue(sqls, argss[:]...)
	}

	batch := db.Pool.SendBatch(context.Background(), &b1)
	for i := 0; i < b1.Len(); i++ {
		exec, err := batch.Exec()
		if err != nil {
			return err
		}
		if exec.RowsAffected() != 1 {
			return errors.New("problem inserting")
		}
	}

	_ = batch.Close()
	return nil
}

func (businesses *businessProfileRepo) PatchProfile(id uuid.UUID, bus *map[string]interface{}) error {
	sql, args, err := sqrl.Update("business_profiles").
		SetMap(*bus).
		Where(sqrl.Eq{"business_id": id}).
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

func (businesses *businessProfileRepo) DeleteProfile(id uuid.UUID) error {
	sql, args, err := sqrl.Delete("business_profiles").
		Where(sqrl.Eq{"business_id": id}).PlaceholderFormat(sqrl.Dollar).ToSql()

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

func (businesses *businessProfileRepo) BuildQuery(p *map[string]interface{}, l *models.Location) (string, []interface{}) {
	keys := make([]string, 0, len(*p))
	tempAssign := make([]string, 0, len(*p))
	values := make([]interface{}, 0, len(*p))
	var i = 1
	for k, v := range *p {
		keys = append(keys, fmt.Sprintf(strings.ToLower(k)))
		tempAssign = append(tempAssign, "$"+strconv.Itoa(i))
		i = i + 1
		values = append(values, v)
	}

	if l != nil {
		keys = append(keys, "location")
		values = append(values, l)
		tempAssign = append(tempAssign, "$"+strconv.Itoa(i))
	}

	sql := fmt.
		Sprintf("INSERT INTO business_profiles (%s) VALUES (%s)", strings.Join(keys, ", "), strings.Join(tempAssign, ","))

	return sql, values
}

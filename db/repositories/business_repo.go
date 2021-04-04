package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/danhper/structomap"
	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
)

var BusinessRepo = &businessRepo{}

type businessRepo struct{}

var columns = []string{"business_name", "business_id", "description", "contact", "location", "picture"}

func (businesses *businessRepo) FindById(id uuid.UUID) (*json_serializer.SmallBusinessResponse, error) {
	tempCol := columns
	tempCol = append(tempCol, "user_id")
	sql, args, err := sqrl.Select(strings.Join(tempCol, ",")).From("businesses_view").Where(sqrl.Eq{"business_id": id}).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	oids := pgx.QueryResultFormatsByOID{db.GeoOid: pgx.BinaryFormatCode}

	args = append(args, 0)
	copy(args[1:], args[0:])
	args[0] = oids

	var business json_serializer.SmallBusinessResponse
	err = pgxscan.Get(context.Background(), db.Pool, &business, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return &business, nil
}

func (businesses *businessRepo) ListBusiness(params *serveUtils.SQLFilter) (*[]interface{}, error) {
	sql, args, err := BuildListQuery(params)
	if err != nil {
		return nil, err
	}

	oids := pgx.QueryResultFormatsByOID{db.GeoOid: pgx.BinaryFormatCode}

	args = append(args, 0)
	copy(args[1:], args[0:])
	args[0] = oids

	if params.Amount == "small" {
		type temp []*json_serializer.SmallBusinessResponse
		var business temp
		if err := pgxscan.Select(context.Background(), db.Pool, &business, sql, args[:]...); err != nil {
			return nil, err
		}

		y := make([]interface{}, len(business))
		for i, v := range business {
			y[i] = v
		}

		return &y, nil
	}
	if params.Amount == "all" {
		var business []*json_serializer.FullBusinessResponse
		if err := pgxscan.Select(context.Background(), db.Pool, &business, sql, args[:]...); err != nil {
			return nil, err
		}

		y := make([]interface{}, len(business))
		for i, v := range business {
			y[i] = v
		}
		return &y, nil
	}

	return nil, errors.New("test")
}

func (businesses *businessRepo) CreateBusiness(business *map[string]interface{}) error {
	sql, args, err := sqrl.Insert("businesses").
		SetMap(*business).PlaceholderFormat(sqrl.Dollar).ToSql()
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

func (businesses *businessRepo) DeleteBusiness(id uuid.UUID) error {
	sql, args, err := sqrl.Delete("businesses").
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

func (businesses *businessRepo) DeleteAddress(id uuid.UUID) error {
	sql, arg, err := sqrl.Delete("addresses").
		Where(sqrl.Eq{"id": id}).PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return err
	}

	exec, err := db.Pool.Exec(context.Background(), sql, arg[:]...)
	if err != nil {
		return err
	}

	count := exec.RowsAffected()
	if count != 1 {
		return err
	}

	return nil
}

func (businesses *businessRepo) BatchDeleteAddress(id []uuid.UUID) error {
	b1 := pgx.Batch{}
	for _, ids := range id {
		sql, arg, err := sqrl.Delete("addresses").
			Where(sqrl.Eq{"id": ids}).PlaceholderFormat(sqrl.Dollar).ToSql()
		if err != nil {
			return err
		}

		b1.Queue(sql, arg[:]...)
	}

	batch := db.Pool.SendBatch(context.Background(), &b1)
	for i := 0; i < b1.Len(); i++ {
		exec, err := batch.Exec()
		if err != nil {
			_ = batch.Close()
			return err
		}
		if exec.RowsAffected() != 1 {
			return errors.New("problem inserting")
		}
	}
	_ = batch.Close()
	return nil
}

func (businesses *businessRepo) PatchBusiness(id uuid.UUID, patch *map[string]interface{}) error {
	// Separating 3 tables, opening hours, owner and address
	var openingHours []models.OpeningHours
	var owner []models.Owner
	var address []models.Address

	if _, ok := (*patch)["opening_hours"]; ok {
		_ = mapstructure.Decode((*patch)["opening_hours"], &openingHours)
		delete(*patch, "opening_hours")
	}
	if _, ok := (*patch)["owner"]; ok {
		_ = mapstructure.Decode((*patch)["owner"], &owner)
		delete(*patch, "owner")
	}
	if _, ok := (*patch)["address"]; ok {
		_ = mapstructure.Decode((*patch)["address"], &address)
		delete(*patch, "address")
	}

	// Separating main business table elements
	temp := map[string]interface{}{}
	if val, ok := (*patch)["approved"]; ok {
		temp["approved"] = val
		delete(*patch, "approved")
	}
	if val, ok := (*patch)["user_id"]; ok {
		temp["user_id"] = val
		delete(*patch, "user_id")
	}
	if val, ok := (*patch)["approved_by"]; ok {
		temp["approved_by"] = val
		delete(*patch, "approved_by")
	}
	if val, ok := (*patch)["blocked"]; ok {
		temp["blocked"] = val
		delete(*patch, "blocked")
	}

	b1 := pgx.Batch{}

	if len(temp) > 0 {
		temp["updated_on"] = time.Now()
		sql, args, err := sqrl.Update("businesses").SetMap(temp).Where(sqrl.Eq{"id": id}).
			PlaceholderFormat(sqrl.Dollar).ToSql()
		if err != nil {
			return err
		}
		b1.Queue(sql, args[:]...)
	}

	if len(openingHours) > 0 {
		for _, val := range openingHours {
			addMap := structomap.New().UseSnakeCase().PickAll().Transform(val)

			sqls, argss, err := sqrl.Insert("opening_times").SetMap(addMap).
				Suffix("ON CONFLICT (business_id,week_day) DO UPDATE SET opening_time = ?, closing_time = ?, opened = ?", val.OpeningTime, val.ClosingTime, val.Opened).
				PlaceholderFormat(sqrl.Dollar).ToSql()
			if err != nil {
				return err
			}
			b1.Queue(sqls, argss[:]...)
		}
	}

	if len(address) > 0 {
		for _, val := range address {
			addMap := structomap.New().UseSnakeCase().PickAll().Transform(val)
			sql, arg, err := sqrl.Insert("addresses").SetMap(addMap).
				Suffix("ON CONFLICT (id) DO UPDATE SET "+
					"business_id = ?, street = ?, ward = ?, municipality = ?, district = ?, state = ?, country = ?, contact = ?",
					val.BusinessId, val.Street, val.Ward, val.Municipality, val.District, val.State, val.Country, val.Contact).
				PlaceholderFormat(sqrl.Dollar).ToSql()
			if err != nil {
				return err
			}
			b1.Queue(sql, arg[:]...)
		}
	}

	if len(*patch) > 0 {
		ab := sqrl.Update("business_profiles").Where(sqrl.Eq{"business_id": id})

		if _, ok := (*patch)["location"]; ok {
			locate := models.Location{}
			_ = mapstructure.Decode((*patch)["location"], &locate)
			delete(*patch, "location")
			(*patch)["location"] = &locate
			ab = ab.SetMap(*patch)
		} else {
			ab = ab.SetMap(*patch)
		}

		sql1, arg1, err1 := ab.PlaceholderFormat(sqrl.Dollar).ToSql()
		if err1 != nil {
			return err1
		}
		b1.Queue(sql1, arg1[:]...)
	}

	batch := db.Pool.SendBatch(context.Background(), &b1)
	for i := 0; i < b1.Len(); i++ {
		exec, err := batch.Exec()
		if err != nil {
			_ = batch.Close()
			return err
		}
		if exec.RowsAffected() != 1 {
			_ = batch.Close()
			return errors.New("problem inserting")
		}
	}

	_ = batch.Close()
	return nil
}

func BuildListQuery(params *serveUtils.SQLFilter) (sql string, args []interface{}, err error) {
	var tempFilter []string
	var tempArgs []interface{}

	for key, val := range params.Filter {
		if key != "distance" && key != "near" {
			tempFilter = append(tempFilter, val)
			tempArgs = append(tempArgs, params.Args[key][0])
		}
	}

	params.Where = append(params.Where, "COUNT(business_id) OVER() AS full_count")

	ab := sqrl.Select(strings.Join(params.Where, ",")).From("businesses_view")

	if _, ok := params.Args["near"]; ok {
		ab := ab.JoinClause("CROSS JOIN (SELECT st_makepoint(?,?)::geography AS ref_geom) AS r",
			params.Args["near"][1], params.Args["near"][0])
		if _, ok := params.Args["distance"]; ok {
			ab = ab.Where("st_dwithin(st_transform(geom, 4326), ref_geom, ?)", params.Args["distance"][0])
		} else {
			ab = ab.Where("st_dwithin(st_transform(geom, 4326), ref_geom, ?)", 5000)
		}
		if val, ok := params.Sort["distance"]; ok {
			ab.OrderBy(val)
			delete(params.Sort, "distance")
		}
	} else {
		if _, ok := params.Sort["distance"]; ok {
			delete(params.Sort, "distance")
		}
	}

	if _, ok := params.Args["approved"]; !ok {
		ab = ab.Where("approved = ?", true)
	}

	var tempSort []string

	for _, val := range params.Sort {
		tempSort = append(tempSort, val)
	}

	if len(tempFilter) > 0 {
		ab = ab.Where(strings.Join(tempFilter, " AND "), tempArgs[:]...)
	}

	if len(tempSort) > 0 {
		ab = ab.OrderBy(strings.Join(tempSort, ","))
	} else {
		ab = ab.OrderBy("created_on DESC")
	}

	ab = ab.Limit(uint64(params.Limit)).Offset(uint64(params.Offset))

	sql, args, err = ab.PlaceholderFormat(sqrl.Dollar).ToSql()
	return
}

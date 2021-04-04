package repositories

import (
	"context"
	"strings"

	"github.com/elgris/sqrl"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
)

var SearchRepo = &searchRepo{}

type searchRepo struct{}

func (user *searchRepo) Search(term string) (interface{}, error) {
	var business []*json_serializer.SmallBusinessResponse

	var value []interface{}
	clause := "approved = ?"
	value = append(value, true)

	clause = clause + " AND category_name LIKE ? OR business_name LIKE ? OR email LIKE ? OR website LIKE ? OR description LIKE ?"
	clause = clause + " OR state LIKE ? OR street LIKE ? OR country LIKE ? OR district LIKE ? OR municipality LIKE ?"

	par := "%" + term + "%"
	value = append(value, par, par, par, par, par, par, par, par, par, par)

	sql, args, err := sqrl.Select(strings.Join(columns, ",")).
		From("businesses_view, jsonb_to_recordset(address) as items(state text, country text, street text, district text, municipality text)").
		Where(clause, value...).
		PlaceholderFormat(sqrl.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	oids := pgx.QueryResultFormatsByOID{db.GeoOid: pgx.BinaryFormatCode}
	args = append(args, 0)
	copy(args[1:], args[0:])
	args[0] = oids

	err = pgxscan.Select(context.Background(), db.Pool, &business, sql, args[:]...)
	if err != nil {
		return nil, err
	}

	return business, nil
}

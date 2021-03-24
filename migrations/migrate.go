package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
)

const version = 7

func Migrate() {
	d := config.Configuration.Database
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", d.DBUser, d.DBPassword, d.DBHost, d.DBPort,
		d.DBName)

	source, err := bindata.WithInstance(bindata.Resource(AssetNames(), Asset))
	if err != nil {
		fmt.Println(err)
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", source, connString)
	if err != nil {
		fmt.Println(err)
	}

	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println(err)
	}
}

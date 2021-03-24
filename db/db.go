package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/log"
	"os"
)

var (
	Pool *pgxpool.Pool
	GeoOid uint32
)

// ConnectDB ...
func init() {
	d := config.Configuration.Database
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", d.DBUser, d.DBPassword, d.DBHost, d.DBPort,
		d.DBName)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Panic()
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = zerologadapter.NewLogger(*log.Logger)

	runtimeParams := make(map[string]string)
	runtimeParams["application_name"] = "pugcha_backend"
	poolConfig.ConnConfig.RuntimeParams = runtimeParams

	Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		os.Exit(1)
	}

	err = Pool.QueryRow(context.Background(), `SELECT 'geo_loc'::regtype::oid`).Scan(&GeoOid)
	if err != nil {
		log.Panic().Err(err)
	}

	go listenChange()
}

func listenChange() {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		log.Err(err)
		return
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN changes")
	if err != nil {
		log.Err(err)
		return
	}
	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			log.Err(err)
			return
		}
		if notification.Payload == "business" {
			cons, err := Pool.Acquire(context.Background())
			if err != nil {
				log.Err(err)
				return
			}
			_, err = cons.Exec(context.Background(), "REFRESH MATERIALIZED VIEW CONCURRENTLY businesses_view")
			if err != nil {
				log.Err(err)
				cons.Release()
			}
			cons.Release()
		} else {
			cons, err := Pool.Acquire(context.Background())
			if err != nil {
				log.Err(err)
			}

			_, err = cons.Exec(context.Background(), "REFRESH MATERIALIZED VIEW CONCURRENTLY users_view")
			if err != nil {
				log.Err(err)
				cons.Release()
			}
			cons.Release()
		}
	}
}

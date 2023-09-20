package database

import (
	"database/sql"
	"fmt"
	"log"

	"go-bunrouter-example/infrastructure/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type HandlerDatabase struct {
	DbConn *bun.DB
}

func NewPostgresDatabaseClient(conf *config.Config) (HandlerDatabase, error) {
	dbConn, err := loadPsqlDb(conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.DBName,
		conf.Env)
	if err != nil {
		log.Printf("failed to connect database instance: %v", err)
		return HandlerDatabase{}, err
	}

	return HandlerDatabase{
		DbConn: dbConn,
	}, nil
}

func loadPsqlDb(hostDB, portDB, userDB, passwordDB, dbName, appName string) (*bun.DB, error) {
	pgConn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", hostDB, portDB)),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithUser(userDB),
		pgdriver.WithPassword(passwordDB),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithApplicationName(appName),
	)

	sqlDB := sql.OpenDB(pgConn)
	db := bun.NewDB(sqlDB, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

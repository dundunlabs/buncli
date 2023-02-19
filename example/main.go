package main

import (
	"database/sql"

	"github.com/dundunlabs/buncli"
	"github.com/dundunlabs/buncli/example/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const dsn = "postgres://dungtd:@localhost:5432/buncli?sslmode=disable"

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	cli := &buncli.Config{
		DB:         db,
		Migrations: migrations.Migrations,
	}
	cli.Run()
}

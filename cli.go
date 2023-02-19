package buncli

import (
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v3"
)

type Config struct {
	*bun.DB
	*migrate.Migrations
}

func (c *Config) Run() {
	migrator := migrate.NewMigrator(c.DB, c.Migrations)
	app := &cli.App{
		Name:  "buncli",
		Usage: "Bun's CLI",
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(ctx *cli.Context) error {
					ms, err := migrator.MigrationsWithStatus(ctx.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "generate",
				Usage: "generate migration",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "sql",
						Usage: "generate sql migration",
					},
				},
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")

					if sql := ctx.Bool("sql"); sql {
						files, err := migrator.CreateSQLMigrations(ctx.Context, name)
						if err != nil {
							return err
						}

						for _, mf := range files {
							fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
						}
					} else {
						mf, err := migrator.CreateGoMigration(ctx.Context, name)
						if err != nil {
							return err
						}
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Migrate database",
				Action: func(ctx *cli.Context) error {
					if err := migrator.Init(ctx.Context); err != nil {
						return err
					}

					if err := migrator.Lock(ctx.Context); err != nil {
						return err
					}

					defer migrator.Unlock(ctx.Context)

					group, err := migrator.Migrate(ctx.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(ctx *cli.Context) error {
					group, err := migrator.Rollback(ctx.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}

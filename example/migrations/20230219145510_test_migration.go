package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type User struct {
	ID   int64 `bun:",pk,autoincrement"`
	Name string
}

type Profile struct {
	ID     int64 `bun:",pk,autoincrement"`
	UserID int64 `bun:",unique,notnull"`
	User   User  `bun:"rel:belongs-to,join:user_id=id"`
}

var models = []any{new(User), new(Profile)}

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		for _, model := range models {
			if _, err := db.NewCreateTable().Model(model).WithForeignKeys().Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		for _, model := range models {
			if _, err := db.NewDropTable().Model(model).Cascade().Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

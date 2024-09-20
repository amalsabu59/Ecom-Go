package migrations

import (
	"context"
	"gologin/internal/models"

	"github.com/uptrace/bun"
)

type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*models.User)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func (m *CreateUsersTable) Down(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().
		Model((*models.User)(nil)).
		IfExists().
		Exec(ctx)
	return err
}

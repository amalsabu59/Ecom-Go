package migrations

import (
	"context"
	"gologin/internal/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type AlterUsersTableIsDeleted struct{}

func (m *AlterUsersTableIsDeleted) Up(ctx context.Context, db *bun.DB) error {
	logger.Log.Info().Msg("Altering users table: Adding is_deleted column")
	_, err := db.ExecContext(ctx, `
        ALTER TABLE users
        ADD COLUMN is_deleted BOOLEAN DEFAULT TRUE;
    `)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to add is_active column")
	}
	return err
}

func (m *AlterUsersTableIsDeleted) Down(ctx context.Context, db *bun.DB) error {
	logger.Log.Info().Msg("Rolling back: Dropping is_deleted column")
	_, err := db.ExecContext(ctx, `
        ALTER TABLE users
        DROP COLUMN IF EXISTS is_deleted;
    `)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to drop is_deleted column")
	}
	return err
}

func init() {

	// // Register the schema change migration.
	alterTableMigrationisDeleted := &AlterUsersTableIsDeleted{}
	Migrations.MustRegister(
		migrate.MigrationFunc(alterTableMigrationisDeleted.Up),
		migrate.MigrationFunc(alterTableMigrationisDeleted.Down),
	)
}

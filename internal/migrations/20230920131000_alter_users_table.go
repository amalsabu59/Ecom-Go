package migrations

import (
	"context"
	"gologin/internal/logger"

	"github.com/uptrace/bun"
)

// AlterUsersTable handles schema changes like adding/removing columns.
type AlterUsersTable struct{}

func (m *AlterUsersTable) Up(ctx context.Context, db *bun.DB) error {
	logger.Log.Info().Msg("Altering users table: Adding is_active column")
	_, err := db.ExecContext(ctx, `
        ALTER TABLE users
        ADD COLUMN is_active BOOLEAN DEFAULT TRUE;
    `)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to add is_active column")
	}
	return err
}

func (m *AlterUsersTable) Down(ctx context.Context, db *bun.DB) error {
	logger.Log.Info().Msg("Rolling back: Dropping is_active column")
	_, err := db.ExecContext(ctx, `
        ALTER TABLE users
        DROP COLUMN IF EXISTS is_active;
    `)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to drop is_active column")
	}
	return err
}

func init() {

	// // Register the schema change migration.
	// alterTableMigration := &AlterUsersTable{}
	// Migrations.MustRegister(
	// 	migrate.MigrationFunc(alterTableMigration.Up),
	// 	migrate.MigrationFunc(alterTableMigration.Down),
	// )
}

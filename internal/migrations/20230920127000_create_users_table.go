package migrations

import (
	"context"
	"gologin/internal/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// CreateUsersTable creates the initial users table.
type CreateUsersTable struct{}

// Up applies the migration to create the users table.
func (m *CreateUsersTable) Up(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*models.User)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

// Down rolls back the migration to drop the users table.
func (m *CreateUsersTable) Down(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().
		Model((*models.User)(nil)).
		IfExists().
		Exec(ctx)
	return err
}

// Register the migrations in the Migrations collection.
func init() {
	// Register the table creation migration.
	createTableMigration := &CreateUsersTable{}
	Migrations.MustRegister(
		migrate.MigrationFunc(createTableMigration.Up),
		migrate.MigrationFunc(createTableMigration.Down),
	)

}

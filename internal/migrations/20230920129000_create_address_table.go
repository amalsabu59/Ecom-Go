package migrations

import (
	"context"
	"gologin/internal/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// CreateUsersTable creates the initial users table.
type CreateAddressTable struct{}

func (m *CreateAddressTable) Up(ctx context.Context, db *bun.DB) error {
	// Create the addresses table with a foreign key reference
	_, err := db.NewCreateTable().
		Model((*models.Address)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *CreateAddressTable) Down(ctx context.Context, db *bun.DB) error {
	// Drop the addresses table
	_, err := db.NewDropTable().Model((*models.Address)(nil)).Exec(ctx)
	return err
}

func init() {
	// Register the table creation migration.
	createAddressMigration := &CreateAddressTable{}
	Migrations.MustRegister(
		migrate.MigrationFunc(createAddressMigration.Up),
		migrate.MigrationFunc(createAddressMigration.Down),
	)

}

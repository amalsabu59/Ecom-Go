package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"gologin/internal/logger"
	"gologin/internal/migrations"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

var DB *bun.DB

// SetupDB initializes the database connection.
func SetupDB() {
	// Configure zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// dsn := os.Getenv("PG_DSN")
	// dsn := "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable"

	dsn := "postgres://myuser:mypassword@host.docker.internal:5432/mydb?sslmode=disable"

	// dsn := "postgresql://amal:WCwxvXnlHqqj5jyodTBuOA@bayou-salmon-5033.8nk.gcp-asia-southeast1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"

	if dsn == "" {
		log.Fatal().Msg("PG_DSN environment variable not set")
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(sqldb, pgdialect.New())

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	logger.Log.Info().Msg("Successfully connected to the database")
}

// SetupTestDB initializes a test database connection.
func SetupTestDB() error {
	// Configure zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dsn := "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(sqldb, pgdialect.New())

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the test database")
	}

	logger.Log.Info().Msg("Successfully connected to the test database")
	return nil
}

// Disconnect closes the database connection.
func Disconnect() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Panic().Err(err).Msg("Failed to close the database connection")
		}
		logger.Log.Info().Msg("Database connection closed successfully")
	}
}

// Migrate applies migrations to the database.
func Migrate() error {
	ctx := context.Background()

	migrator := migrate.NewMigrator(DB, migrations.Migrations)

	// Initialize migrations
	if err := migrator.Init(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize migrations")
	}

	// Apply migrations
	migrated, err := migrator.Migrate(ctx) // Corrected variable name
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Migration failed")
	}

	logger.Log.Info().Msgf("Migrations applied successfully: %v", migrated)
	return nil
}

// Rollback reverts the last migration.
func Rollback() error {
	ctx := context.Background()

	migrator := migrate.NewMigrator(DB, migrations.Migrations)

	// Rollback the last migration
	rolledBack, err := migrator.Rollback(ctx)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Rollback failed")
	}

	logger.Log.Info().Msgf("Rollbacks applied successfully: %v", rolledBack)
	return nil
}

// CreateSQLMigration generates SQL migration files for PostgreSQL.
func CreateSQLMigration(name string) error {
	ctx := context.Background()

	migrator := migrate.NewMigrator(DB, migrations.Migrations)

	// Generate migration file name based on input arguments
	migrationName := strings.Join(strings.Fields(name), "_")

	// Create SQL migrations for both "up" and "down"
	files, err := migrator.CreateSQLMigrations(ctx, migrationName)
	if err != nil {
		return fmt.Errorf("failed to create SQL migration: %w", err)
	}

	// Log the created migration files
	for _, file := range files {
		logger.Log.Info().Msgf("Created migration: %s (%s)", file.Name, file.Path)
	}

	return nil
}

// main function to handle command-line arguments

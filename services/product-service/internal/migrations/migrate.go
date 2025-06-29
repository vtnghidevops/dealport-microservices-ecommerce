package migrations

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// NewMigrator tạo một migrate instance mới
func NewMigrator(db *sqlx.DB) (*migrate.Migrate, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	// Lấy đường dẫn migration từ biến môi trường hoặc sử dụng mặc định
	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		// Mặc định là thư mục migrations ngang hàng với thư mục cmd
		migrationPath = "file://migrations"
	}

	log.Printf("Using migration path: %s", migrationPath)

	// Tạo driver cho postgres
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Tạo migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// RunMigrations chạy tất cả migrations lên version mới nhất
func RunMigrations(db *sqlx.DB) error {
	start := time.Now()
	log.Println("Starting database migration...")

	m, err := NewMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Chạy migration
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Lấy version hiện tại
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	elapsed := time.Since(start)
	if err == migrate.ErrNoChange {
		log.Printf("No migration needed, current version: %d, took %s", version, elapsed)
	} else {
		log.Printf("Migration completed successfully to version %d (dirty: %t), took %s", version, dirty, elapsed)
	}

	return nil
}

// RollbackMigration roll back lại 1 version
func RollbackMigration(db *sqlx.DB) error {
	log.Println("Rolling back last database migration...")

	m, err := NewMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Roll back 1 version
	err = m.Steps(-1)
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	// Lấy version hiện tại
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Rollback completed successfully to version %d (dirty: %t)", version, dirty)
	return nil
}

// DropAllTables xóa tất cả các bảng - chỉ dùng cho development
func DropAllTables(db *sqlx.DB) error {
	if os.Getenv("ENVIRONMENT") != "development" {
		return errors.New("dropping tables is only allowed in development environment")
	}

	log.Println("Dropping all tables...")

	m, err := NewMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Drop tất cả
	err = m.Drop()
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// GetMigrationStatus trả về trạng thái migration hiện tại
func GetMigrationStatus(db *sqlx.DB) (uint, bool, error) {
	m, err := NewMigrator(db)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}

	version, dirty, err := m.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// ForceVersion cưỡng chế đặt version migration - dùng để fix dirty database
func ForceVersion(db *sqlx.DB, version int) error {
	log.Printf("Forcing migration version to %d...", version)

	m, err := NewMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Force version
	err = m.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("Successfully forced migration version to %d", version)
	return nil
}

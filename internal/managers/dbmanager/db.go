package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cash-flow-financial/internal/models"

	_ "github.com/lib/pq"
)

type DBManager struct {
	db     *sql.DB
	config *models.DatabaseConfig
}

func NewDBManager(cfg *models.DatabaseConfig) (*DBManager, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection to %s:%s: %w", cfg.Host, cfg.Port, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection established successfully to %s:%s", cfg.Host, cfg.Port)

	manager := &DBManager{
		db:     db,
		config: cfg,
	}

	return manager, nil
}

func (dm *DBManager) GetDB() *sql.DB {
	return dm.db
}

func (dm *DBManager) Close() error {
	if dm.db != nil {
		log.Println("Closing database connection...")
		return dm.db.Close()
	}
	return nil
}

func (dm *DBManager) IsHealthy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := dm.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (dm *DBManager) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := dm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (dm *DBManager) ExecuteQuery(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := dm.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return result, nil
}

func (dm *DBManager) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return dm.db.QueryRowContext(ctx, query, args...)
}

func (dm *DBManager) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dm.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return rows, nil
}

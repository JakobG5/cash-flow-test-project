package dbmanager

import (
	"context"
	"database/sql"
)

type IDBManager interface {
	GetDB() *sql.DB
	Close() error
	IsHealthy() error
	WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error
}

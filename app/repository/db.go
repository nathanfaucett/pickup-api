package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aicacia/pickup/app/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB() error {
	connection, err := sqlx.Connect("postgres", env.GetDatabaseUrl())
	if err != nil {
		return err
	}
	db = connection
	return nil
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	err := db.Close()
	if err != nil {
		return err
	}
	db = nil
	return nil
}

func StringFromSQLNullString(nullString sql.NullString) *string {
	if nullString.Valid {
		return &nullString.String
	}
	return nil
}

func ValidConnection() bool {
	return db.Ping() == nil
}

func TimeFromSQLNullTime(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

func All[T any](query string, args ...interface{}) ([]T, error) {
	var rows []T
	err := db.Select(&rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetOptional[T any](query string, args ...interface{}) (*T, error) {
	rows, err := All[T](query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		return &rows[0], nil
	}
	return nil, nil
}

func Get[T any](query string, args ...interface{}) (T, error) {
	var row T
	err := db.Get(&row, query, args...)
	if err != nil {
		return row, err
	}
	return row, nil
}

func Transaction[T any](fn func(tx *sqlx.Tx) (T, error)) (T, error) {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return *new(T), txErr
	}
	result, fnErr := fn(tx)
	if fnErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return *new(T), errors.Join(fnErr, rollbackErr)
		}
		return *new(T), fnErr
	} else {
		commitErr := tx.Commit()
		if commitErr != nil {
			return *new(T), commitErr
		}
		return result, nil
	}
}

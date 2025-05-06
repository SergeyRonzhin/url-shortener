package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	db     *sqlx.DB
	logger *logger.Logger
}

var (
	queryTimeout = 5 * time.Second
	tableName    = "urls"
	initQuery    = `CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		short_url VARCHAR(255) NOT NULL,
		original_url VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
)

func NewDBStorage(options *config.Options, logger *logger.Logger) (*DBStorage, error) {

	db, err := sqlx.Connect("postgres", options.DatabaseDsn)

	if err != nil {
		return nil, err
	}

	if err = db.DB.Ping(); err != nil {
		return nil, err
	}

	db.MustExec(initQuery)

	return &DBStorage{
		db:     db,
		logger: logger,
	}, nil
}

func (s *DBStorage) Add(url URL) error {

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := s.db.NamedExecContext(ctx, `INSERT INTO `+tableName+` (uuid, short_url, original_url) 
		VALUES (:uuid, :short_url, :original_url)`, url)

	return err
}

func (s *DBStorage) Batch(urls []URL) error {

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := s.db.NamedExecContext(ctx, `INSERT INTO `+tableName+` (uuid, short_url, original_url)
        VALUES (:uuid, :short_url, :original_url)`, urls)

	return err
}

func (s *DBStorage) GetShortURL(original string) (bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.db.QueryRowContext(ctx, "SELECT short_url FROM "+tableName+" WHERE original_url = $1", original)
	var short string

	err := row.Scan(&short)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, ""
	case err != nil:
		s.logger.Error(err)
		return false, ""
	default:
		return true, short
	}
}

func (s *DBStorage) GetOriginalURL(short string) (bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM "+tableName+" WHERE short_url = $1", short)
	var original string

	err := row.Scan(&original)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, ""
	case err != nil:
		s.logger.Error(err)
		return false, ""
	default:
		return true, original
	}
}

func (s DBStorage) Close() error {
	return s.db.Close()
}

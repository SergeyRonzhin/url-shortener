package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	db      *sqlx.DB
	options *config.Options
	logger  *logger.Logger
}

var (
	initQuery = `CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	
	CREATE TABLE IF NOT EXISTS urls (
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
		db:      db,
		options: options,
		logger:  logger,
	}, nil
}

func (s *DBStorage) Add(ctx context.Context, url URL) error {

	_, err := s.db.NamedExecContext(ctx, `INSERT INTO urls (uuid, short_url, original_url) 
		VALUES (:uuid, :short_url, :original_url)`, url)

	return err
}

func (s *DBStorage) Batch(ctx context.Context, urls []URL) error {

	_, err := s.db.NamedExecContext(ctx, `INSERT INTO urls (uuid, short_url, original_url)
        VALUES (:uuid, :short_url, :original_url)`, urls)

	return err
}

func (s *DBStorage) GetShortURL(ctx context.Context, original string) (bool, string) {
	return s.getUrl(ctx, "SELECT short_url FROM urls WHERE original_url = $1", original)
}

func (s *DBStorage) GetOriginalURL(ctx context.Context, short string) (bool, string) {
	return s.getUrl(ctx, "SELECT original_url FROM urls WHERE short_url = $1", short)
}

func (s DBStorage) Close() error {
	return s.db.Close()
}

func (s *DBStorage) Ping(ctx context.Context) error {

	db, err := sqlx.Connect("postgres", s.options.DatabaseDsn)

	defer func() {
		err = db.Close()
	}()

	if err != nil {
		return err
	}

	err = db.DB.PingContext(ctx)

	return err
}

func (s *DBStorage) getUrl(ctx context.Context, query string, args ...any) (bool, string) {
	row := s.db.QueryRowContext(ctx, query, args...)
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

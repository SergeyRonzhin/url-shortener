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

func NewDBStorage(options *config.Options, logger *logger.Logger) (*DBStorage, error) {

	db, err := sqlx.Connect("postgres", options.DatabaseDsn)

	if err != nil {
		return nil, err
	}

	if err = db.DB.Ping(); err != nil {
		return nil, err
	}

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
	return s.getURL(ctx, "SELECT short_url FROM urls WHERE original_url = $1", original)
}

func (s *DBStorage) GetOriginalURL(ctx context.Context, short string) (bool, string) {
	return s.getURL(ctx, "SELECT original_url FROM urls WHERE short_url = $1", short)
}

func (s DBStorage) Close() error {
	return s.db.Close()
}

func (s *DBStorage) Ping(ctx context.Context) error {

	db, err := sqlx.ConnectContext(ctx, "postgres", s.options.DatabaseDsn)

	defer func() {
		err = db.Close()
	}()

	return err
}

func (s *DBStorage) getURL(ctx context.Context, query string, args ...any) (bool, string) {
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

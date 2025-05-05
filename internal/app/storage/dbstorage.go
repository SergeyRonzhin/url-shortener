package storage

import (
	"context"
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
	tableName = "urls"
	initQuery = `CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	
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

	err = initDatabase(db)

	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db:     db,
		logger: logger,
	}, nil
}

func (s *DBStorage) Get(key string) (string, bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM "+tableName+" WHERE short_url = $1", key)
	var url string

	err := row.Scan(&url)

	if err != nil {
		s.logger.Error(err)
		return "", false
	}

	return url, true
}

func (s *DBStorage) Add(key string, value string) error {

	_, err := s.db.Exec("INSERT INTO "+tableName+" (short_url, original_url) VALUES ($1, $2)", key, value)
	return err
}

func (s *DBStorage) ContainsValue(value string) (bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, "SELECT short_url FROM "+tableName+" WHERE original_url = $1", value)
	var url string

	err := row.Scan(&url)

	if err != nil {
		s.logger.Error(err)
		return false, ""
	}

	return true, url
}

func (s DBStorage) Close() error {
	return s.db.Close()
}

func initDatabase(db *sqlx.DB) error {

	_, err := db.Exec(initQuery)
	return err
}

package storage

import (
	"context"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	DB     *sqlx.DB
	Logger *logger.Logger
}

type DBURL struct {
	UUID        uuid.UUID `db:"uuid"`
	ShortURL    string    `db:"short_url"`
	OriginalURL string    `db:"original_url"`
	CreatedAt   time.Time `db:"created_at"`
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
		DB:     db,
		Logger: logger,
	}, nil
}

func (s *DBStorage) Get(key string) (string, bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, "SELECT original_url FROM "+tableName+" WHERE short_url = $1", key)
	var url string

	err := row.Scan(&url)

	if err != nil {
		s.Logger.Error(err)
		return "", false
	}

	return url, true
}

func (s *DBStorage) Add(key string, value string) error {

	_, err := s.DB.Exec("INSERT INTO "+tableName+" (short_url, original_url) VALUES ($1, $2)", key, value)
	return err
}

func (s *DBStorage) ContainsValue(value string) (bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, "SELECT short_url FROM "+tableName+" WHERE original_url = $1", value)
	var url string

	err := row.Scan(&url)

	if err != nil {
		s.Logger.Error(err)
		return false, ""
	}

	return true, url
}

func initDatabase(db *sqlx.DB) error {

	_, err := db.Exec(initQuery)
	return err
}

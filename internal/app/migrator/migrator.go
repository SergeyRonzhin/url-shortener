package migrator

import (
	"errors"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	options *config.Options
	logger  *logger.Logger
}

func New(options *config.Options, logger *logger.Logger) *Migrator {
	return &Migrator{options, logger}
}

func (m Migrator) ApplyMigrations() error {

	if m.options.MigrationsPath == "" {
		return errors.New("migrations path is required")
	}

	fm, err := migrate.New("file://"+m.options.MigrationsPath, m.options.DatabaseDsn)

	if err != nil {
		return err
	}

	if err := fm.Up(); err != nil {

		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to apply")
			return nil
		}

		return err
	}

	m.logger.Info("migrations applied successfully")
	return nil
}

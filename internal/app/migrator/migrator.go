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

	defer func() {
		if fm != nil {
			sourceErr, databaseErr := fm.Close()

			if sourceErr != nil || databaseErr != nil {
				err = errors.Join(sourceErr, databaseErr)
			}
		}
	}()

	if err != nil {
		return err
	}

	err = fm.Up()

	if err == nil {
		m.logger.Info("migrations applied successfully")
	} else if errors.Is(err, migrate.ErrNoChange) {

		m.logger.Info("no migrations to apply")
		return nil
	}

	return err
}

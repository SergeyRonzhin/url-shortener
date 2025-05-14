package migrator

import (
	"errors"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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

	sourceDriver, err := iofs.New(migrations.MigrationsFiles, m.options.MigrationsPath)

	if err != nil {
		return err
	}

	fm, err := migrate.NewWithSourceInstance("iofs", sourceDriver, m.options.DatabaseDsn)

	defer func() {
		if fm == nil {
			return
		}

		sourceErr, databaseErr := fm.Close()

		if sourceErr != nil || databaseErr != nil {
			err = errors.Join(err, sourceErr, databaseErr)
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

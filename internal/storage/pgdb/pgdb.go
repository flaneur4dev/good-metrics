package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
	// "github.com/flaneur4dev/good-metrics/internal/lib/utils"
)

type DBStorage struct {
	db  *sql.DB
	key string
}

func New(source, k string) (*DBStorage, error) {
	db, err := sql.Open("pgx", source)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	ds := &DBStorage{
		db:  db,
		key: k,
	}

	return ds, nil
}

func (ds *DBStorage) AllMetrics() (gm, cm []string) {
	return
}

func (ds *DBStorage) OneMetric(t, n string) (cs.Metrics, error) {
	return cs.Metrics{}, e.ErrNoUsedDB
}

func (ds *DBStorage) Update(n string, nm cs.Metrics) (cs.Metrics, error) {
	return cs.Metrics{}, e.ErrNoUsedDB
}

func (ds *DBStorage) Check() error {
	if err := ds.db.Ping(); err != nil {
		return fmt.Errorf("can't connect to database: %w", err)
	}
	return nil
}

func (ds *DBStorage) Close() error {
	if err := ds.db.Close(); err != nil {
		return fmt.Errorf("can't close database: %w", err)
	}
	return nil
}

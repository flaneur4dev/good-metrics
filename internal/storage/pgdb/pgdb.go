package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
)

type DBStorage struct {
	db  *sql.DB
	key string
}

func New(source, k string) (*DBStorage, error) {
	db, err := sql.Open("pgx", source)
	if err != nil {
		return nil, fmt.Errorf(e.NoDBOpen, err)
	}

	// if err = db.Ping(); err != nil {
	// 	return nil, fmt.Errorf(e.NoDBConnect, err)
	// }

	ds := &DBStorage{
		db:  db,
		key: k,
	}

	if err = ds.initDB(); err != nil {
		return nil, fmt.Errorf(e.NoDBInit, err)
	}

	return ds, nil
}

func (ds *DBStorage) AllMetrics() (gm, cm []string) {
	rows, err := ds.db.Query("SELECT id, type, value, delta FROM metrics")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			i, t string
			v    sql.NullFloat64
			d    sql.NullInt64
		)

		if err := rows.Scan(&i, &t, &v, &d); err != nil {
			fmt.Println(err)
			return
		}

		if t == utils.GaugeName {
			gm = append(gm, fmt.Sprintf("%s: %f", i, v.Float64))
		} else {
			cm = append(cm, fmt.Sprintf("%s: %d", i, d.Int64))
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (ds *DBStorage) OneMetric(t, n string) (cs.Metrics, error) {
	var (
		m = cs.Metrics{}
		v sql.NullFloat64
		d sql.NullInt64
	)

	err := ds.db.QueryRow("SELECT * FROM metrics WHERE id=$1 AND type=$2 LIMIT 1", n, t).Scan(&m.ID, &m.MType, &v, &d, &m.Hash)
	switch {
	case err == sql.ErrNoRows:
		return cs.Metrics{}, e.ErrNoMetric
	case err != nil:
		return cs.Metrics{}, fmt.Errorf(e.NoMetricFetch, err)
	}

	switch m.MType {
	case utils.GaugeName:
		val := cs.Gauge(v.Float64)
		m.Value = &val
	case utils.CounterName:
		del := cs.Counter(d.Int64)
		m.Delta = &del
	}

	return m, nil
}

func (ds *DBStorage) Update(nm cs.Metrics) (cs.Metrics, error) {
	if err := utils.ValidateMetric(nm, ds.key); err != nil {
		return cs.Metrics{}, err
	}

	switch nm.MType {
	case utils.CounterName:
		return ds.updateCounter(nm)
	default:
		return ds.updateGauge(nm)
	}
}

func (ds *DBStorage) Check() error {
	if err := ds.db.Ping(); err != nil {
		return fmt.Errorf(e.NoDBConnect, err)
	}
	return nil
}

func (ds *DBStorage) Close() error {
	if err := ds.db.Close(); err != nil {
		return fmt.Errorf(e.NoDBClose, err)
	}
	return nil
}

func (ds *DBStorage) initDB() error {
	q := `CREATE TABLE IF NOT EXISTS metrics (
		id varchar(50) PRIMARY KEY,
		type varchar(10) NOT NULL,
		value float8 NULL,
		delta int NULL,
		hash char(64)
	);`

	_, err := ds.db.Exec(q)
	if err != nil {
		return fmt.Errorf(e.NoDBTable, err)
	}
	return nil
}

func (ds *DBStorage) updateGauge(m cs.Metrics) (cs.Metrics, error) {
	result, err := ds.db.Exec("UPDATE metrics SET value=$1, hash=$2 WHERE id=$3 AND type=$4", *m.Value, m.Hash, m.ID, m.MType)
	if err != nil {
		return cs.Metrics{}, fmt.Errorf(e.NoMetricUpdate, m.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return cs.Metrics{}, fmt.Errorf(e.NoMetricUpdate, m.ID, err)
	}

	if rows == 0 {
		res, err := ds.db.Exec(
			"INSERT INTO metrics (id, type, value, hash) VALUES ($1, $2, $3, $4)",
			m.ID, m.MType, *m.Value, m.Hash,
		)
		if err != nil {
			return cs.Metrics{}, fmt.Errorf(e.NoMetricInsert, m.ID, err)
		}

		if rr, err := res.RowsAffected(); err != nil || rr != 1 {
			return cs.Metrics{}, fmt.Errorf(e.NoMetricInsert, m.ID, err)
		}
	}

	return m, nil
}

func (ds *DBStorage) updateCounter(m cs.Metrics) (cs.Metrics, error) {
	var d sql.NullInt64

	err := ds.db.QueryRow("SELECT delta FROM metrics WHERE id=$1 AND type=$2 LIMIT 1", m.ID, m.MType).Scan(&d)
	switch {
	case err == sql.ErrNoRows:
		res, err := ds.db.Exec(
			"INSERT INTO metrics (id, type, delta, hash) VALUES ($1, $2, $3, $4)",
			m.ID, m.MType, *m.Delta, m.Hash,
		)
		if err != nil {
			return cs.Metrics{}, fmt.Errorf(e.NoMetricInsert, m.ID, err)
		}

		if rs, err := res.RowsAffected(); err != nil || rs != 1 {
			return cs.Metrics{}, fmt.Errorf(e.NoMetricInsert, m.ID, err)
		}

		return m, nil
	case err != nil:
		return cs.Metrics{}, fmt.Errorf(e.NoMetricUpdate, m.ID, err)
	}

	nv := cs.Counter(d.Int64) + *m.Delta
	m.Delta = &nv

	msg := fmt.Sprintf(utils.CounterTmpl, m.ID, *m.Delta)
	m.Hash = utils.Sign256(msg, ds.key)

	result, err := ds.db.Exec("UPDATE metrics SET delta=$1, hash=$2 WHERE id=$3 AND type=$4", *m.Delta, m.Hash, m.ID, m.MType)
	if err != nil {
		return cs.Metrics{}, fmt.Errorf(e.NoMetricUpdate, m.ID, err)
	}

	if rows, err := result.RowsAffected(); err != nil || rows != 1 {
		return cs.Metrics{}, fmt.Errorf(e.NoMetricUpdate, m.ID, err)
	}

	return m, nil
}

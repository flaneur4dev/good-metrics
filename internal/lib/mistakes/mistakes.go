package mistakes

import "errors"

const (
	NoDBOpen    = "can't open database: %w"
	NoDBConnect = "can't connect to database: %w"
	NoDBInit    = "can't init database: %w"
	NoDBClose   = "can't close database: %w"
	NoDBTable   = "can't create table: %w"

	NoMetricFetch  = "can't fetch metric: %w"
	NoMetricInsert = "can't insert %s metric: %w"
	NoMetricUpdate = "can't update %s metric: %w"
)

var (
	ErrNoMetric         = errors.New("no such metric")
	ErrUnkownMetricType = errors.New("unknown metric type")
	ErrInvalidData      = errors.New("invalid data")
	ErrCompromisedData  = errors.New("compromised data")
	ErrNoUsedDB         = errors.New("database not used")
)

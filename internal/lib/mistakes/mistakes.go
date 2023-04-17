package mistakes

import "errors"

var (
	ErrNoMetric         = errors.New("no such metric")
	ErrUnkownMetricType = errors.New("unknown metric type")
	ErrInvalidData      = errors.New("invalid data")
	ErrCompromisedData  = errors.New("compromised data")
	ErrNoUsedDB         = errors.New("database not used")
)

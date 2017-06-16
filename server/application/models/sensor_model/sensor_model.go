package sensor_model

import (
	"database/sql"
)

type Table struct {
	Id int
	Name string
	ControllerId sql.NullInt64
	ActivationDate sql.NullString
	Status int
	DeactivationDate sql.NullString
	SensorType int
	Company string
}

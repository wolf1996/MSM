package sensor

import (
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type RegisterSensorForm struct {
	ControllerId int64 `json:"controller_id"`
	SensorId     int64 `json:"sensor_id"`
}

func (form *RegisterSensorForm) Validate() int {
	if form.ControllerId == 0 || form.SensorId == 0 {
		return error_codes.INVALID_JSON
	}
	return error_codes.OK
}

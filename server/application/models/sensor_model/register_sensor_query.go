package sensor_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

func RegisterSensorQuery(controllerId, sensorId int64) models.ErrorModel {
	qr, err := models.Database.Query("UPDATE SENSOR "+
		"SET controller_id = $1, status = 1, activation_date = NOW() "+
		"WHERE id = $2 RETURNING id", controllerId, sensorId)
	if err != nil {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	if !qr.Next() {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Error no sensor or controller with such id", err), Code: error_codes.DATABASE_INVALID_SENSOR}
	}
	return nil
}

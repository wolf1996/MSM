package sensor_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"fmt"
)

func RegisterSensorQuery(controllerId, sensorId, userId int64) models.ErrorModel {
	checkQuery := "SELECT id FROM controllers WHERE user_id = $1"
	qr, err := models.Database.Query(checkQuery, userId)
	defer qr.Close()
	if err != nil {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	var expectedId int64
	for qr.Next() {
		if err = qr.Scan(&expectedId); err != nil {
			return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
		}
	}
	if controllerId != expectedId {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Ownership error ", err), Code: error_codes.INVALID_OWNER}
	}
	updateQuery := "UPDATE SENSOR SET controller_id = $1, status = 1, activation_date = NOW() WHERE id = $2 RETURNING id"
	qr, err = models.Database.Query(updateQuery, controllerId, sensorId)
	defer qr.Close()
	if err != nil {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	if !qr.Next() {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Error no sensor or controller with such id", err), Code: error_codes.DATABASE_INVALID_SENSOR}
	}
	return nil
}

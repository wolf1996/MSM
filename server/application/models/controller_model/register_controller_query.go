package controller_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

func RegisterControllerQuery(userId, controllerId int64) models.ErrorModel {
	qr, err := models.Database.Query("UPDATE CONTROLLERS "+
		"SET user_id = $1, status = 1, activation_date = NOW() "+
		"WHERE id = $2 RETURNING id", userId, controllerId)
	if err != nil {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	if !qr.Next() {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Error no controller with such id", err), Code: error_codes.DATABASE_INVALID_CONTROLLER}
	}
	return nil
}

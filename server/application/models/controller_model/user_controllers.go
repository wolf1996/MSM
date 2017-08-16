package controller_model

import (
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type ControllerModel Table

type ControllerModels []ControllerModel

func GetUserControllers(id int64) (ControllerModels, models.ErrorModel) {
	var infoSlice ControllerModels
	qr, err := models.Database.Query(
		"SELECT * "+
			"FROM CONTROLLERS WHERE user_id = $1 ;", id)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info ControllerModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.UserId, &info.Adres, &info.ActivationDate,
			&info.Status, &info.Mac, &info.DeactivationDate, &info.ControllerType)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}

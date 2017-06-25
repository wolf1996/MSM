package controller_model

import(
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
)

type ControllerModel Table

type ControllerModels []ControllerModel

func GetUserControllers( id int64 ) (ControllerModels, models.ErrorModel){
	var infoSlice ControllerModels
	qr,err := models.Database.Query(
		"SELECT * " +
			"FROM CONTROLLERS WHERE user_id = $1 ;", id)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	defer qr.Close()
	var info ControllerModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.UserId, &info.Adres, &info.ActivationDate,
			&info.Status, &info.Mac, &info.DeactivationDate, &info.ControllerType)
		if err != nil{
			return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}
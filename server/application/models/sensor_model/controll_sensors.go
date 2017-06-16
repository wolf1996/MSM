package sensor_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
)

type SensorModel Table
type SensorModels []SensorModel

func GetControlledSensors( controllerId, userId int64 ) (SensorModels, models.ErrorModel){
	var infoSlice SensorModels
	qr,err := models.Database.Query(
		" SELECT "+
	" SENSOR.id, SENSOR.Name, SENSOR.controller_id, SENSOR.activation_date, SENSOR.status, SENSOR.deactivation_date, SENSOR.sensor_type, SENSOR.company"+
	" FROM SENSOR INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id WHERE controller_id = $1 AND user_id = $2;", controllerId, userId)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	defer qr.Close()
	var info SensorModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.ControllerId, &info.ActivationDate,
			&info.Status, &info.DeactivationDate, &info.SensorType,  &info.Company)
		if err != nil{
			return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}
package data_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
)

type DataModel Table

type DataModels []DataModel

func GetData(userId, sensorId int64, date  string, limit int64 ) (DataModels, models.ErrorModel) {
	var infoSlice DataModels
	qr,err := models.Database.Query(" SELECT"+
	" DATA.sensor_id, DATA.date, DATA.value, DATA.hs"+
	" FROM DATA INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id"+
	" INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id"+
	" WHERE sensor_id = $1 AND user_id = $2 AND date < $3 "+
	" LIMIT $4", sensorId, userId, date, limit)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	defer qr.Close()
	var info DataModel
	for qr.Next() {
		err = qr.Scan(&info.SensorId, &info.Date, &info.Value, &info.Hash)
		if err != nil{
			return infoSlice, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}

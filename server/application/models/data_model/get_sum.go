package data_model

import (
	"database/sql"
	"fmt"
	"MSM/server/application/models"
	"MSM/server/application/error_codes"
)

type DeltaData struct {
	Delta sql.NullFloat64
}

func GetSum(userId, sensorId int64, dateBegin string, dateEnd string) (DeltaData, models.ErrorModel) {
	var info DeltaData
	qr, err := models.Database.Query(" SELECT"+
		" (max(DATA.value)-min(DATA.value))"+
		" FROM DATA INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id"+
		" INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id"+
		" WHERE sensor_id = $1 AND user_id = $2 AND date >= $3 AND date  < $4 "+
		" ;", sensorId, userId, dateBegin, dateEnd)
	if err != nil {
		return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	for qr.Next() {
		err = qr.Scan(&info.Delta)
		if err != nil {
			return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
	}
	return info, nil
}

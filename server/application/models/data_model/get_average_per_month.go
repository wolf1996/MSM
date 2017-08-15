package data_model

import (
	"database/sql"
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"MSM/server/application/error_codes"
)

type AveragePerMonth struct {
	Average sql.NullFloat64
}

func GetAveragePerMonth(userId, sensorId int64, dateBegin string, dateEnd string) (AveragePerMonth, models.ErrorModel) {
	var info AveragePerMonth
	qr, err := models.Database.Query(" WITH mnths AS (SELECT extract(MONTH FROM date) AS mnth, "+
		"(max(DATA.value)-min(DATA.value)) AS value "+
		"FROM DATA "+
		"INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id "+
		"INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id "+
		"WHERE sensor_id = $1 AND user_id = $2 AND date >= $3 "+
		"                 AND date < $4 "+
		"           GROUP BY 1 "+
		") "+
		"SELECT avg(value) "+
		"FROM mnths", sensorId, userId, dateBegin, dateEnd)
	if err != nil {
		return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	for qr.Next() {
		err = qr.Scan(&info.Average)
		if err != nil {
			return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
	}
	return info, nil
}

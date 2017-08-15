package data_model

import (
	"database/sql"
	"fmt"
	"MSM/server/application/models"
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
		return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: 2}
	}
	defer qr.Close()
	for qr.Next() {
		err = qr.Scan(&info.Average)
		if err != nil {
			return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: 2}
		}
	}
	return info, nil
}

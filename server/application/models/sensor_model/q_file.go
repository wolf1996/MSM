package sensor_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/framework"
	"fmt"
)

type SensorQueries struct{ *models.MainDatabase}


func GetSensorQueries(cont framework.AppContext) *SensorQueries{
	db  := cont.GetValue(models.DbSystemName)
		if db == nil {
		panic(fmt.Sprintf("Can't find %s in resources", db))
	}
	database, ok := db.(models.MainDatabase)
		if !ok {
		panic(fmt.Sprintf("Can't convert %s, wrong type", db))
	}
	return &SensorQueries{&database}
}
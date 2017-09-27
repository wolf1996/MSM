package controller_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/framework"
	"fmt"
)

type ControllerQueries struct {*models.MainDatabase}

func GetControllerQueries(cont framework.AppContext) *ControllerQueries{
	db  := cont.GetValue(models.DbSystemName)
	if db == nil {
		panic(fmt.Sprintf("Can't find %s in resources", db))
	}
	database, ok := db.(models.MainDatabase)
	if !ok {
		panic(fmt.Sprintf("Can't convert %s, wrong type", db))
	}
	return &ControllerQueries{&database}
}
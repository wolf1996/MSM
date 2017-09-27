package user_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
	"github.com/wolf1996/MSM/server/framework"
)

type UserQueries struct {*models.MainDatabase}

func GetUserQueries(cont framework.AppContext) *UserQueries {
	db  := cont.GetValue(models.DbSystemName)
	if db == nil {
		panic(fmt.Sprintf("Can't find %s in resources", db))
	}
	database, ok := db.(models.MainDatabase)
	if !ok {
		panic(fmt.Sprintf("Can't convert %s, wrong type", db))
	}
	return &UserQueries{&database,}
}

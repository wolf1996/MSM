package application

import (
	_ "github.com/wolf1996/MSM/server/application/controllers"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/framework"
	"net/http"
)

func prepareResource(port, dbLogin, dbPass, dbURL string)(cnt framework.AppContext, err error) {
	db,err :=  models.GetDatabase(dbLogin, dbPass, dbURL)
	if err != nil {
		return
	}
	res := framework.Resource{
		models.DbSystemName,
		db,
	}
	framework.AddResource(res)
	return framework.GetContext()
}

func AppStart(port, dbLogin, dbPass, dbURL string) error{
	var cnt framework.AppContext
	cnt, err := prepareResource(port, dbLogin, dbPass, dbURL)
	if err != nil {
		return  err
	}
	router := framework.HandlerConstructor(cnt)
	http.ListenAndServe(port, router)
	return nil
}

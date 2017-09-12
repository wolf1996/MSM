package application

import (
	_ "github.com/wolf1996/MSM/server/application/controllers"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/framework"
	"net/http"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	var cnt framework.AppContext
	router := framework.HandlerConstructor(cnt)
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

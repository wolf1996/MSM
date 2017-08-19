package application

import (
	"github.com/wolf1996/MSM/server/application/models"
	"net/http"
	"github.com/wolf1996/MSM/server/framework"
	_ "github.com/wolf1996/MSM/server/application/controllers"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	router := framework.GetRouters()
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

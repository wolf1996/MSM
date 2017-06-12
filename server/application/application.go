package application

import (
	"github.com/wolf1996/MSM/server/application/controllers"
	"net/http"
	"github.com/wolf1996/MSM/server/application/models"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	router := controllers.GetRouters()
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

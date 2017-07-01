package application

import (
	"github.com/wolf1996/MSM/server/application/controllers"
	"github.com/wolf1996/MSM/server/application/models"
	"net/http"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	router := controllers.GetRouters()
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

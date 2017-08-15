package application

import (
	"MSM/server/application/controllers"
	"MSM/server/application/models"
	"net/http"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	router := controllers.GetRouters()
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

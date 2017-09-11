package application

import (
	_ "github.com/wolf1996/MSM/server/application/controllers"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/framework"
	"net/http"
	"context"
)

func AppStart(port, dbLogin, dbPass, dbURL string) {
	var cnt context.Context
	router := framework.HandlerConstructor(cnt)
	models.Init(dbLogin, dbPass, dbURL)
	http.ListenAndServe(port, router)
}

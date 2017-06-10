package application

import (
	"github.com/wolf1996/MSM/server/application/controllers"
	"log"
	"net/http"
)

func AppStart(port string) {
	router := controllers.GetRouters()
	log.Fatal(http.ListenAndServe(port, router))
}

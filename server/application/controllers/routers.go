package controllers

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routTable Routes

func GetRouters() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routTable {
		logsystem.Info.Printf("%s registered to path %s", route.Name, route.Pattern)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func AddRout(rt Route) {
	routTable = append(routTable, rt)
}

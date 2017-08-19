package framework

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
)

type HandlerFunc http.HandlerFunc

type MiddleWare func(HandlerFunc) HandlerFunc

type Route struct {
	Name        string
	Method      string
	Pattern     string
	MidleWare   []MiddleWare
	HandlerFunc HandlerFunc
}

type Routes []Route

var routTable Routes

func applyMiddlewares(handlerFunc HandlerFunc, middlewares []MiddleWare) HandlerFunc{
	for _, i := range middlewares {
		handlerFunc = i(handlerFunc)
	}
	return handlerFunc
}

func GetRouters() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routTable {
		handler := applyMiddlewares(route.HandlerFunc, route.MidleWare)
		logsystem.Info.Printf("%s registered to path %s", route.Name, route.Pattern)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(http.HandlerFunc(handler))
	}
	return router
}

func AddRout(rt Route) {
	routTable = append(routTable, rt)
}

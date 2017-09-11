package framework

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"context"
)

type lowHandlerFunc http.HandlerFunc

type ContHandlerFunc func (appcontext context.Context, w http.ResponseWriter, r *http.Request)

type MiddleWare func(ContHandlerFunc) ContHandlerFunc

type Route struct {
	Name        string
	Method      string
	Pattern     string
	MidleWare   []MiddleWare
	HandlerFunc ContHandlerFunc
}

type Routes []Route

var routTable Routes

func applyMiddlewares(handlerFunc ContHandlerFunc, middlewares []MiddleWare) lowHandlerFunc{
	for _, i := range middlewares {
		handlerFunc = i(handlerFunc)
	}
	var cont context.Context
	return  AppContextMiddleware(cont,handlerFunc)
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

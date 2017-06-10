package controllers

import (
	"fmt"
	"net/http"
)

func init() {
	rout := Route{"TestController", "GET", "/controller/test", testController}
	AddRout(rout)
}

func testController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

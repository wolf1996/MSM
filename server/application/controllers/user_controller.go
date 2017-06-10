package controllers

import (
	"fmt"
	"net/http"
)

func init() {
	rout := Route{"TestUser", "GET", "/user/test", test}
	AddRout(rout)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

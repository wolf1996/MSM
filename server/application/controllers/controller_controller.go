package controllers

import (
	"fmt"
	"net/http"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/logsystem"
	"github.com/wolf1996/MSM/server/application/models/controller_model"
	"github.com/wolf1996/MSM/server/application/view/controller"
)

func init() {
	rout := Route{"TestController", "GET", "/controller/test", testController}
	AddRout(rout)
	rout = Route{"ControllersInfo", "GET", "/controller/get_user_controllers", getUserController}
	AddRout(rout)
}

func testController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func getUserController(w http.ResponseWriter, r *http.Request){
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w,nil, 2)
		sess,_ = session_manager.NewSession(r,"user_session")
		sess.Save(r,w)
		return
	}
	id := sess.Values["user"]
	if id == nil{
		logsystem.Error.Printf("LogIn first")
		view.WriteMessage(&w,view.ErrorMsg{"Login first"}, 1)
		return
	}
	md, errDb := controller_model.GetUserControllers(id.(int64))
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		view.WriteMessage(&w,view.ErrorMsg{"Database Error"}, 2)
		return
	}
	var inf []controller.ControllerInfo
	for _,i := range md{
		buf := controller.ControllerInfo{i.Id, i.Name, i.UserId, i.Adres,
			i.ActivationDate.String, i.Status, i.Mac, i.DeactivationDate.String,
			i.ControllerType}
		inf = append(inf, buf)
	}
	view.WriteMessage(&w,inf, 0)
}
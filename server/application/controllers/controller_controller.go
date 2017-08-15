package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"MSM/server/application/models/controller_model"
	"MSM/server/application/models/sensor_model"
	"MSM/server/application/session_manager"
	"MSM/server/application/view"
	"MSM/server/application/view/controller"
	"MSM/server/logsystem"
	"net/http"
	"strconv"
)

func init() {
	rout := Route{"TestController", "GET", "/controller/test", testController}
	AddRout(rout)
	rout = Route{"ControllersInfo", "GET", "/controller/get_user_controllers", getUserController}
	AddRout(rout)
	rout = Route{"ControllersInfo", "GET", "/controller/{id}/get_controller_stats", getControllerView}
	AddRout(rout)
}

func testController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func compileControllerInfo(model *controller_model.ControllerModel) *controller.ControllerInfo{
	var activationDate *string
	var deactivationDate *string
	if model.ActivationDate.Valid{
		activationDate = &model.ActivationDate.String
	}
	if model.DeactivationDate.Valid{
		deactivationDate = &model.DeactivationDate.String
	}
	return &controller.ControllerInfo{&model.Id,
							&model.Name,
							&model.UserId,
							&model.Adres,
							activationDate,
							&model.Status,
							&model.Mac,
							deactivationDate,
							&model.ControllerType,
	}
}

func getUserController(w http.ResponseWriter, r *http.Request) {
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Session Error"}, 2)
		sess, _ = session_manager.NewSession(r, "user_session")
		w.WriteHeader(http.StatusForbidden)
		sess.Save(r, w)
		return
	}
	id := sess.Values["user"]
	if id == nil {
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login first"}, 1)
		return
	}
	md, errDb := controller_model.GetUserControllers(id.(int64))
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, 2)
		return
	}
	var inf []controller.ControllerInfo
	for _, i := range md {
		buf := compileControllerInfo(&i)
		inf = append(inf, *buf)
	}
	view.WriteMessage(&w, inf, 0)
}

func compileControllerStats(id *int64, month, prevMonth, prevYear *float64) *controller.ControllerStats{
	return &controller.ControllerStats{id,
		                              month,
		                              prevMonth,
		                              prevYear}
}


func getControllerView(w http.ResponseWriter, r *http.Request) {
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Session Error"}, 2)
		w.WriteHeader(http.StatusForbidden)
		sess, _ = session_manager.NewSession(r, "user_session")
		sess.Save(r, w)
		return
	}
	id, ok := sess.Values["user"].(int64)
	if !ok {
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login first"}, 1)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	cId := vals["id"]
	controllerId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	sensors, errDb := sensor_model.GetTaxedSensors(controllerId, id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, 2)
		return
	}

	var month, year, prevMonth float64
	for _, i := range sensors {
		stats, errCd := getSensorStats(id, i.Id)
		if errCd != nil {
			switch {
			case errCd.Id() == 1:
				logsystem.Error.Printf("Database Error %s", errCd.Error())
				w.WriteHeader(http.StatusInternalServerError)
				view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, errCd.Id())
				return

			}
		}
		if stats.CurrentMonth != nil {
			month += i.Tax * *stats.CurrentMonth
		}
		if stats.PrevYear != nil {
			year += i.Tax * *stats.PrevYear
		}
		if stats.PrevYearMonth != nil {
			prevMonth += i.Tax * *stats.PrevYearMonth
		}
	}
	vw := compileControllerStats(&controllerId, &month, &prevMonth, &year)
	view.WriteMessage(&w, vw, 0)
}

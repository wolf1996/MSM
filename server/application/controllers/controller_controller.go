package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/models/controller_model"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/controller"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"encoding/json"
	"io/ioutil"
	"io"
)

func init() {
	rout := Route{"TestController", "GET", "/controller/test", testController}
	AddRout(rout)
	rout = Route{"ControllersInfo", "GET", "/controller/get_user_controllers", getUserController}
	AddRout(rout)
	rout = Route{"ControllersInfo", "GET", "/controller/{id}/get_controller_stats", getControllerView}
	AddRout(rout)
	rout = Route{"RegisterController", "POST", "/controller/register", registerController}
	AddRout(rout)
}

func testController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func compileControllerInfo(model controller_model.ControllerModel) controller.ControllerInfo{
	var activationDate *string
	var deactivationDate *string
	if model.ActivationDate.Valid{
		activationDate = &model.ActivationDate.String
	}
	if model.DeactivationDate.Valid{
		deactivationDate = &model.DeactivationDate.String
	}
	return controller.ControllerInfo{&model.Id,
							&model.Name,
							&model.ObjectId,
							&model.Meta,
							activationDate,
							&model.Status,
							&model.Mac,
							deactivationDate,
							&model.ControllerType,
	}
}

func registerController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logsystem.Error.Printf("Post Json loading in registerController %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Body Read"}, error_codes.INVALID_BODY_READ)
		return
	}
	form := new(controller.RegisterControllerForm)
	if err = json.Unmarshal(body, form); err != nil {
		logsystem.Error.Printf("Unmarshal error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Unmarshal error"}, error_codes.UNMARSHAL_ERROR)
		return
	}
	if form.ControllerId == 0 {
		logsystem.Error.Printf("Invalid fields in json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Invalid fields in json %s"}, error_codes.INVALID_JSON)
		return
	}
	session, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Session Error"}, error_codes.DATABASE_ERROR)
		session, _ = session_manager.NewSession(r, "user_session")
		w.WriteHeader(http.StatusForbidden)
		session.Save(r, w)
		return
	}
	userId := session.Values["user"]
	if userId == nil {
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login first"}, error_codes.NOT_LOGGED)
		return
	}
	// logsystem.Info.Printf("%v", form)
	if err_code := controller_model.RegisterControllerQuery(userId.(int64), form.ControllerId, form.ObjectId); err_code != nil {
		logsystem.Error.Printf("Controller registration failed %s", err_code)
		w.WriteHeader(http.StatusConflict)
		view.WriteMessage(&w, view.ErrorMsg{"Controller Registration Failed"}, error_codes.LOGIN_FAILED)
		return
	}
	view.WriteMessage(&w, view.ErrorMsg{"Controller Registration Success!"}, error_codes.OK)
}

func getUserController(w http.ResponseWriter, r *http.Request) {
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Session Error"}, error_codes.DATABASE_ERROR)
		sess, _ = session_manager.NewSession(r, "user_session")
		w.WriteHeader(http.StatusForbidden)
		sess.Save(r, w)
		return
	}
	id := sess.Values["user"]
	if id == nil {
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login first"}, error_codes.NOT_LOGGED)
		return
	}
	md, errDb := controller_model.GetUserControllers(id.(int64))
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var inf []controller.ControllerInfo
	for _, i := range md {
		buf := compileControllerInfo(i)
		inf = append(inf, buf)
	}
	view.WriteMessage(&w, inf, error_codes.OK)
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
		view.WriteMessage(&w, view.ErrorMsg{"Session Error"}, error_codes.SESSION_ERROR)
		w.WriteHeader(http.StatusForbidden)
		sess, _ = session_manager.NewSession(r, "user_session")
		sess.Save(r, w)
		return
	}
	id, ok := sess.Values["user"].(int64)
	if !ok {
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login first"}, error_codes.NOT_LOGGED)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	cId := vals["id"]
	controllerId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	sensors, errDb := sensor_model.GetTaxedSensors(controllerId, id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}

	var month, year, prevMonth float64
	for _, i := range sensors {
		stats, errCd := getSensorStats(id, i.Id)
		if errCd != nil {
			switch {
			case errCd.Id() == error_codes.DATABASE_ERROR:
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
	view.WriteMessage(&w, vw, error_codes.OK)
}

package controllers

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/models/object_model"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/object"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/view/controller"
	"github.com/wolf1996/MSM/server/application/models/controller_model"
	"io/ioutil"
	"io"
	"encoding/json"
)

func init() {
	rout := Route{"getUserObjects", "GET", "/object/get_user_objects", getUserObjects}
	AddRout(rout)
	rout = Route{"getObjectControllers", "GET", "/object/{id}/get_object_controllers", getObjectControllers}
	AddRout(rout)
	rout = Route{"getObjectView", "GET", "/object/{id}/get_object_stats", getObjectView}
	AddRout(rout)
	rout = Route{"RegisterObject", "POST", "/object/register", registerObject}
	AddRout(rout)
}
func compileObjectInfo(model object_model.ObjectModel) object.ObjectInfo{
	return object.ObjectInfo{&model.Id,
							&model.Name,
							&model.UserId,
							&model.Addres,

	}
}

func getUserObjects(w http.ResponseWriter, r *http.Request) {
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
	md, errDb := object_model.GetUserObjects(id.(int64))
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var inf []object.ObjectInfo
	for _, i := range md {
		buf := compileObjectInfo(i)
		inf = append(inf, buf)
	}
	view.WriteMessage(&w, inf, error_codes.OK)
}

func compileObjectControllers(id *int64, month, prevMonth, prevYear *float64) *controller.ControllerStats{
	return &controller.ControllerStats{id,
		                              month,
		                              prevMonth,
		                              prevYear}
}


func getObjectView(w http.ResponseWriter, r *http.Request) {
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

func getObjectControllers(w http.ResponseWriter, r *http.Request) {
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
	objectId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	controllers, errDb := controller_model.GetObjectControllers(id, objectId)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var controllersInfo controller.ControllersInfo
	for _, v := range controllers {
		inf := compileControllerInfo(v)
		controllersInfo = append(controllersInfo, inf)
	}
	view.WriteMessage(&w, controllersInfo, error_codes.OK)
}


func registerObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logsystem.Error.Printf("Post Json loading in registerController %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Body Read"}, error_codes.INVALID_BODY_READ)
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

	form := new(object.RegisterForm)
	if err = json.Unmarshal(body, form); err != nil {
		logsystem.Error.Printf("Unmarshal error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Unmarshal error"}, error_codes.UNMARSHAL_ERROR)
		return
	}

	if err_code := object_model.RegisterObjectQuery(userId.(int64), form.Name , form.Adres); err_code != nil {
		logsystem.Error.Printf("Object registration failed %s", err_code)
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Object Registration Failed"}, error_codes.LOGIN_FAILED)
		return
	}
	view.WriteMessage(&w, view.ErrorMsg{"Object Registration Success!"}, error_codes.OK)
}
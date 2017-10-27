package controllers

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/sensor"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"io/ioutil"
	"io"
	"encoding/json"
)

func init() {
	rout := Route{"ControllersInfo", "GET", "/controller/{id}/get_sensors", getControllerSensor}
	AddRout(rout)
	rout = Route{"RegisterController", "POST", "/sensor/register", registerSensor}
	AddRout(rout)
}

func compileSensorInfo(v *sensor_model.SensorModel) *sensor.SensorInfo {
	var deactivationDate,activationDate *string
	var controllerId *int64
	if v.ActivationDate.Valid{
		activationDate = &v.ActivationDate.String
	}
	if v.DeactivationDate.Valid{
		deactivationDate = &v.DeactivationDate.String
	}
	if v.ControllerId.Valid {
		controllerId = &v.ControllerId.Int64
	}
	return &sensor.SensorInfo{&v.Id,
					  &v.Name,
		              controllerId,
		              activationDate,
		              &v.Status,
					  deactivationDate,
		              &v.SensorType,
					  &v.Company,
	}
}

func registerSensor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logsystem.Error.Printf("Post Json loading in registerController %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Body Read"}, error_codes.INVALID_BODY_READ)
		return
	}
	form := new(sensor.RegisterSensorForm)
	if err = json.Unmarshal(body, form); err != nil {
		logsystem.Error.Printf("Unmarshal error %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Unmarshal error"}, error_codes.UNMARSHAL_ERROR)
		return
	}
	if form.Validate() != error_codes.OK {
		logsystem.Error.Printf("Invalid fields in json %s", err)
		w.WriteHeader(http.StatusBadRequest)
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
	if err = sensor_model.RegisterSensorQuery(form.ControllerId, form.SensorId, userId.(int64)); err != nil {
		logsystem.Error.Printf("Controller registration failed %s", err)
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Controller Registration Failed"}, error_codes.LOGIN_FAILED)
		return
	}
	view.WriteMessage(&w, view.ErrorMsg{"Controller Registration Success!"}, error_codes.OK)
}

func getControllerSensor(w http.ResponseWriter, r *http.Request) {
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
	sensors, errDb := sensor_model.GetControlledSensors(controllerId, id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var sensorsInfo sensor.SensorsInfo
	for _, v := range sensors {
		inf := *compileSensorInfo(&v)
		sensorsInfo = append(sensorsInfo, inf)
	}
	view.WriteMessage(&w, sensorsInfo, error_codes.OK)
}

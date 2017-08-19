package controllers

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/sensor"
	"github.com/wolf1996/MSM/server/framework"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
)

func init() {
	rout := framework.Route{Name: "ControllersInfo",
		Method:      "GET",
		Pattern:     "/controller/{id}/get_sensors",
		HandlerFunc: getControllerSensor,
	}
	framework.AddRout(rout)
}

func compileSensorInfo(v *sensor_model.SensorModel) *sensor.SensorInfo {
	var deactivationDate, activationDate *string
	var controllerId *int64
	if v.ActivationDate.Valid {
		activationDate = &v.ActivationDate.String
	}
	if v.DeactivationDate.Valid {
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

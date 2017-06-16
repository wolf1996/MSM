package controllers

import (
	"net/http"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/logsystem"
	"strconv"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/view/sensor"
	"github.com/gorilla/mux"
)

func init() {
	rout := Route{"ControllersInfo", "GET", "/controller/{id}/get_sensors", getControllerSensor}
	AddRout(rout)
}


//TODO СДЕЛАТЬ ЛИСТ СЕНСОРОВ

func getControllerSensor(w http.ResponseWriter, r *http.Request){
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w,nil, 2)
		sess,_ = session_manager.NewSession(r,"user_session")
		sess.Save(r,w)
		return
	}
	id, ok := sess.Values["user"].(int64)
	if !ok{
		logsystem.Error.Printf("LogIn first")
		view.WriteMessage(&w,view.ErrorMsg{"Login first"}, 1)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		view.WriteMessage(&w,view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	cId := vals["id"]
	controllerId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		view.WriteMessage(&w,view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	sensors, errDb := sensor_model.GetControlledSensors(controllerId, id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		view.WriteMessage(&w,view.ErrorMsg{"Database Error"}, 2)
		return
	}
	var sensorsInfo sensor.SensorsInfo
	for _,v := range sensors {
		inf := sensor.SensorInfo{v.Id, v.Name, v.ControllerId.Int64, v.ActivationDate.String, v.Status,
		v.DeactivationDate.String, v.SensorType, v.Company}
		sensorsInfo = append(sensorsInfo, inf)
	}
	view.WriteMessage(&w,sensorsInfo, 0)
}

package controllers

import (
	"net/http"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/logsystem"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/models/data_model"
	"time"
	"github.com/wolf1996/MSM/server/application/view/data"
)

func init() {
	rout := Route{"ControllersInfo", "GET", "/sensor/{id}/get_data", getSensorData}
	AddRout(rout)
}


func getSensorData(w http.ResponseWriter, r *http.Request){
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w,view.ErrorMsg{"Session Error"}, 2)
		w.WriteHeader(http.StatusForbidden)
		sess,_ = session_manager.NewSession(r,"user_session")
		sess.Save(r,w)
		return
	}
	id, ok := sess.Values["user"].(int64)
	if !ok{
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w,view.ErrorMsg{"Login first"}, 1)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w,view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	cId := vals["id"]
	sensorId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		view.WriteMessage(&w,view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}
	query := r.URL.Query()
	limitSt := query["limit"]
	if len(limitSt) > 1 {
		logsystem.Error.Printf("invalid query")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w,view.ErrorMsg{"invalid query"}, 4)
		return
	}
	if len(limitSt) == 0{
		limitSt = append(limitSt, "100")
	}
	limit, err := strconv.ParseInt(limitSt[0], 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w,view.ErrorMsg{"Can't parse argument"}, 3)
		return
	}

	dateSt := query["date"]
	if len(dateSt) > 1 {
		logsystem.Error.Printf("invalid query")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w,view.ErrorMsg{"invalid query"}, 4)
		return
	}

	if len(dateSt) == 0{
		nowtime := time.Now().Local()
		strtime := nowtime.Format("2006-01-02")
		dateSt = append(dateSt, strtime)
	}
	date := dateSt[0]
	dataLst, errDb := data_model.GetData(id, sensorId, date,limit)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w,view.ErrorMsg{"Database Error"}, 2)
		return
	}
	var dataInfoList data.DataInfoList
	for _,v := range dataLst {
		inf := data.DataInfo{v.SensorId, v.Date, v.Value, v.Hash}
		dataInfoList = append(dataInfoList, inf)
	}
	view.WriteMessage(&w, dataInfoList, 0)
}

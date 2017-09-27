package controllers

import (
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/models/data_model"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/data"
	"github.com/wolf1996/MSM/server/framework"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
	"time"
)

func init() {
	rout := framework.Route{Name: "view stats",
		Method:      "GET",
		Pattern:     "/sensor/{id}/view_stats",
		HandlerFunc: getSensorStatsData,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
	rout = framework.Route{Name: "get data",
		Method:      "GET",
		Pattern:     "/sensor/{id}/get_data",
		HandlerFunc: getSensorData,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
}

func compileDataInfoStats(month, prev *data_model.DeltaData, year *data_model.AveragePerMonth) *data.DataInfoStats {
	var monthReal, prevReal, yearReal *float64
	if month.Delta.Valid {
		monthReal = &month.Delta.Float64
	}
	if prev.Delta.Valid {
		prevReal = &prev.Delta.Float64
	}
	if year.Average.Valid {
		yearReal = &year.Average.Float64
	}
	return &data.DataInfoStats{monthReal, prevReal, yearReal}
}

func getSensorStats(appContext framework.AppContext,id int64, sensorId int64) (data.DataInfoStats, models.ErrorModel) {
	nowtime := time.Now()
	firstDay := time.Date(nowtime.Year(), nowtime.Month(), 1, 0, 0, 0, 0, nowtime.Location())
	firstNextDay := nowtime.AddDate(0, 1, 0)
	strBeg := firstDay.Format("2006-01-02")
	strEnd := firstNextDay.Format("2006-01-02")
	sumPerMonth, errDb := data_model.GetDataQueries(appContext).GetSum(id, sensorId, strBeg, strEnd)
	if errDb != nil {
		return data.DataInfoStats{}, models.ErrorModelImpl{errDb.Error(), error_codes.DATABASE_ERROR}
	}
	firstDay = firstDay.AddDate(-1, 0, 0)
	firstNextDay = firstDay.AddDate(0, 1, 0)
	strBeg = firstDay.Format("2006-01-02")
	strEnd = firstNextDay.Format("2006-01-02")
	sumPerPrevMonth, errDb := data_model.GetDataQueries(appContext).GetSum(id, sensorId, strBeg, strEnd)
	if errDb != nil {
		return data.DataInfoStats{}, models.ErrorModelImpl{errDb.Error(), error_codes.DATABASE_ERROR}
	}
	currYear := time.Date(nowtime.Year(), 1, 1, 0, 0, 0, 0, nowtime.Location())
	firstPrevYear := currYear.AddDate(-1, 0, 0)
	strBeg = firstPrevYear.Format("2006-01-02")
	strEnd = currYear.Format("2006-01-02")
	yearAverPerMonth, errDb := data_model.GetDataQueries(appContext).GetAveragePerMonth(id, sensorId, strBeg, strEnd)
	if errDb != nil {
		return data.DataInfoStats{}, models.ErrorModelImpl{errDb.Error(), error_codes.DATABASE_ERROR}
	}
	stats := compileDataInfoStats(&sumPerMonth, &sumPerPrevMonth, &yearAverPerMonth)
	return *stats, nil
}

func getSensorStatsData(appContext framework.AppContext,w http.ResponseWriter, r *http.Request) {
	cont := r.Context()
	id,ok := cont.Value("id").(int64)
	if !ok {
		logsystem.Error.Printf("No Id in context")
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Server Error"}, error_codes.SERVER_ERROR)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	cId := vals["id"]
	sensorId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}

	stats, errCd := getSensorStats(appContext, id, sensorId)
	if errCd != nil {
		switch {
		case errCd.Id() == 1:
			logsystem.Error.Printf("Database Error %s", errCd.Error())
			w.WriteHeader(http.StatusInternalServerError)
			view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
			return

		}
	}
	sensorInfo, errCd := sensor_model.GetSensorQueries(appContext).GetTaxedSensor(sensorId, id)
	if errCd != nil {
		logsystem.Error.Printf("Invalid sensor")
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Invalid sensor"}, error_codes.DATABASE_INVALID_SENSOR)
		return
	}
	var accural float32
	if stats.CurrentMonth != nil {
		accural = float32((sensorInfo.Tax) * (*stats.CurrentMonth))
	} else {
		accural = 0
	}
	overpay := float32(10.0)
	rl := float32(accural - overpay)
	info := compileSensorVidgetData(sensorInfo, accural, overpay, rl, stats)
	view.WriteMessage(&w, *info, error_codes.OK)
}

func compileSensorVidgetData(model sensor_model.SensorTaxedModel,
	accural, overpay, rl float32, stats data.DataInfoStats) *data.SensorVidgetData {
	result := &data.SensorVidgetData{&model.Type,
		&model.Name,
		&model.Status,
		&accural,
		&overpay,
		&rl,
		&stats}
	return result
}

func getSensorData(appContext framework.AppContext, w http.ResponseWriter, r *http.Request) {
	cont := r.Context()
	id,ok := cont.Value("id").(int64)
	if !ok {
		logsystem.Error.Printf("No Id in context")
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Server Error"}, error_codes.SERVER_ERROR)
		return
	}
	vals := mux.Vars(r)
	if vals == nil {
		logsystem.Error.Printf("Can't parse argument")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	cId := vals["id"]
	sensorId, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}
	query := r.URL.Query()
	limitSt := query["limit"]
	if len(limitSt) > 1 {
		logsystem.Error.Printf("invalid query")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"invalid query"}, error_codes.INVALID_QUERY)
		return
	}
	if len(limitSt) == 0 {
		limitSt = append(limitSt, "100")
	}
	limit, err := strconv.ParseInt(limitSt[0], 10, 64)
	if err != nil {
		logsystem.Error.Printf("Can't parse argument %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Can't parse argument"}, error_codes.INVALID_ARGUMENT)
		return
	}

	dateSt := query["date"]
	if len(dateSt) > 1 {
		logsystem.Error.Printf("invalid query")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"invalid query"}, error_codes.INVALID_QUERY)
		return
	}

	if len(dateSt) == 0 {
		nowtime := time.Now().Local()
		strtime := nowtime.Format("2006-01-02")
		dateSt = append(dateSt, strtime)
	}
	date := dateSt[0]
	dataLst, errDb := data_model.GetDataQueries(appContext).GetData(id, sensorId, date, limit)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var dataInfoList data.DataInfoList
	for _, v := range dataLst {
		inf := compileView(v)
		dataInfoList = append(dataInfoList, *inf)
	}
	view.WriteMessage(&w, dataInfoList, error_codes.OK)
}

func compileView(model data_model.DataModel) (result *data.DataInfo) {
	result = &data.DataInfo{&model.SensorId,
		&model.Date,
		&model.Value,
		&model.Hash}
	return
}

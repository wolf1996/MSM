package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models/controller_model"
	"github.com/wolf1996/MSM/server/application/models/sensor_model"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/controller"
	"github.com/wolf1996/MSM/server/framework"
	"github.com/wolf1996/MSM/server/logsystem"
	"net/http"
	"strconv"
)

func init() {
	rout := framework.Route{Name: "TestController",
		Method:      "GET",
		Pattern:     "/controller/test",
		HandlerFunc: testController,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
	rout = framework.Route{Name: "ControllersInfo",
		Method:      "GET",
		Pattern:     "/controller/get_user_controllers",
		HandlerFunc: getUserController,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
	rout = framework.Route{Name: "ControllersInfo",
		Method:      "GET",
		Pattern:     "/controller/{id}/get_controller_stats",
		HandlerFunc: getControllerView,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
}

func testController(appContext framework.AppContext,w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func compileControllerInfo(appContext framework.AppContext,model *controller_model.ControllerModel) *controller.ControllerInfo {
	var activationDate *string
	var deactivationDate *string
	if model.ActivationDate.Valid {
		activationDate = &model.ActivationDate.String
	}
	if model.DeactivationDate.Valid {
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

func getUserController(appContext framework.AppContext,w http.ResponseWriter, r *http.Request) {
	cont := r.Context()
	id,ok := cont.Value("id").(int64)
	if !ok {
		logsystem.Error.Printf("No Id in context")
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Server Error"}, error_codes.SERVER_ERROR)
		return
	}
	md, errDb := controller_model.GetUserControllers(id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.DATABASE_ERROR)
		return
	}
	var inf []controller.ControllerInfo
	for _, i := range md {
		buf := compileControllerInfo(appContext,&i)
		inf = append(inf, *buf)
	}
	view.WriteMessage(&w, inf, error_codes.OK)
}

func compileControllerStats(id *int64, month, prevMonth, prevYear *float64) *controller.ControllerStats {
	return &controller.ControllerStats{id,
		month,
		prevMonth,
		prevYear}
}

func getControllerView(appContext framework.AppContext,w http.ResponseWriter, r *http.Request) {
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

package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/models/user_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/view/user"
	"github.com/wolf1996/MSM/server/framework"
	"github.com/wolf1996/MSM/server/logsystem"
	"io"
	"io/ioutil"
	"net/http"
	"context"
)

func init() {
	rout := framework.Route{Name: "TestUser",
		Method:      "GET",
		Pattern:     "/user/test",
		HandlerFunc: test,
	}
	framework.AddRout(rout)
	rout = framework.Route{Name: "SignInUser",
		Method:      "POST",
		Pattern:     "/user/sign_in",
		HandlerFunc: signIn,
	}
	framework.AddRout(rout)
	rout = framework.Route{Name: "GetUserInfo",
		Method:      "GET",
		Pattern:     "/user/user_info",
		HandlerFunc: getUserInfo,
		MidleWare: []framework.MiddleWare{framework.AuthRequired,},
	}
	framework.AddRout(rout)
}

func test(appContext context.Context,w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func signIn(appContext context.Context,w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logsystem.Error.Printf("Post Json loading in signIn %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Body Read"}, error_codes.INVALID_BODY_READ)
		return

	}
	if err := r.Body.Close(); err != nil {
		logsystem.Error.Printf("Body Close %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Body Close Error"}, error_codes.INVALID_BODY_CLOSE)
		return
	}
	logIn := user.LoginForm{}
	if err := json.Unmarshal(body, &logIn); err != nil {
		logsystem.Error.Printf("Unmarshal error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Unmarshal error"}, error_codes.UNMARSHAL_ERROR)
		return
	}
	isValid := logIn.Validate()
	if !isValid {
		logsystem.Error.Printf("Invalid")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Validation Failed"}, error_codes.VALIDATION_FAILED)
		return
	}
	var id int64
	var loggerr models.ErrorModel
	if id, loggerr = user_model.LogInUser(*logIn.EMail, *logIn.Pass); loggerr != nil {
		logsystem.Error.Printf("Login failed %s", loggerr)
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login Failed"}, error_codes.LOGIN_FAILED)
		return
	}
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Session error"}, error_codes.SESSION_ERROR)
		sess.Save(r, w)
		return
	}
	sess.Values["user"] = id
	sess.Save(r, w)
	view.WriteMessage(&w, view.ErrorMsg{"Ok"}, error_codes.OK)
}

func compileUserInfo(info *user_model.UserInfoModel) *user.UserInfo {
	var familyName, name, secondName, dateReceiving, issuedBy *string
	var divisionNumber, registrationAddres, mailingAddres, birthday *string
	var sex *bool
	var homePhone, mobilePhone, citizenShip *string
	if info.FamilyName.Valid {
		familyName = &info.FamilyName.String
	}
	if info.Name.Valid {
		name = &info.Name.String
	}
	if info.SecondName.Valid {
		secondName = &info.SecondName.String
	}
	if info.IssuedBy.Valid {
		issuedBy = &info.IssuedBy.String
	}
	if info.DivisionNumber.Valid {
		divisionNumber = &info.DivisionNumber.String
	}
	if info.RegistrationAddres.Valid {
		registrationAddres = &info.RegistrationAddres.String
	}
	if info.MailingAddres.Valid {
		mailingAddres = &info.MailingAddres.String
	}
	if info.BirthDay.Valid {
		birthday = &info.BirthDay.String
	}
	if info.Sex.Valid {
		sex = &info.Sex.Bool
	}
	if info.HomePhone.Valid {
		homePhone = &info.HomePhone.String
	}
	if info.MobilePhone.Valid {
		mobilePhone = &info.MobilePhone.String
	}
	if info.CitizenShip.Valid {
		citizenShip = &info.CitizenShip.String
	}
	return &user.UserInfo{familyName,
		name,
		secondName,
		dateReceiving,
		issuedBy,
		divisionNumber,
		registrationAddres,
		mailingAddres,
		birthday,
		sex,
		homePhone,
		mobilePhone,
		citizenShip,
		&info.EMail,
	}
}

func getUserInfo(appContext context.Context,w http.ResponseWriter, r *http.Request) {
	cont := r.Context()
	id,ok := cont.Value("id").(int64)
	if !ok {
		logsystem.Error.Printf("No Id in context")
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Server Error"}, error_codes.SERVER_ERROR)
		return
	}
	md, errDb := user_model.UserInfoQuery(id)
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w, view.ErrorMsg{"Database Error"}, error_codes.SERVER_ERROR)
		return
	}
	inf := compileUserInfo(&md)
	view.WriteMessage(&w, inf, error_codes.OK)
}

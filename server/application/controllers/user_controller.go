package controllers

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"github.com/wolf1996/MSM/server/application/view/user"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/models/user_model"
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/logsystem"
)

func init() {
	rout := Route{"TestUser", "GET", "/user_model/test", test}
	AddRout(rout)
	rout = Route{"SignInUser", "POST", "/user_model/sign_in", signIn}
	AddRout(rout)
	rout = Route{"GetUserInfo", "GET", "/user_model/user_info", getUserInfo}
	AddRout(rout)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body,1048576))
	if err != nil{
		logsystem.Error.Printf("Post Json loading in signIn %s", err)
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w,view.ErrorMsg{"Body Read"}, 1)
		return

	}
	if err := r.Body.Close(); err != nil {
		logsystem.Error.Printf("Body Close %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w,view.ErrorMsg{"Body Close Error"}, 2)
		return
	}
	logIn := user.LoginForm{}
	if err := json.Unmarshal(body, &logIn); err!= nil{
		logsystem.Error.Printf("Unmarshal error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w,view.ErrorMsg{"Unmarshal error"}, 3)
		return
	}
	isValid := logIn.Validate()
	if ! isValid{
		logsystem.Error.Printf("Invalid")
		w.WriteHeader(http.StatusBadRequest)
		view.WriteMessage(&w, view.ErrorMsg{"Validation Failed"}, 4)
		return
	}
	var id int64
	var loggerr models.ErrorModel
	if id, loggerr = user_model.LogInUser(logIn.EMail, logIn.Pass); loggerr != nil{
		logsystem.Error.Printf("Login failed %s", loggerr)
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w, view.ErrorMsg{"Login Failed"}, 5)
		return
	}
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w,view.ErrorMsg{"Session error"}, 6)
		sess.Save(r,w)
		return
	}
	sess.Values["user"] = id
	sess.Save(r,w)
	view.WriteMessage(&w,view.ErrorMsg{"Ok"}, 0)
}

func getUserInfo(w http.ResponseWriter, r *http.Request)()  {
	sess, err := session_manager.GetSession(r, "user_session")
	if err != nil {
		logsystem.Error.Printf("Get session error %s", err)
		view.WriteMessage(&w,nil, 2)
		sess,_ = session_manager.NewSession(r,"user_session")
		w.WriteHeader(http.StatusForbidden)
		sess.Save(r,w)
		return
	}
	id := sess.Values["user"]
	if id == nil{
		logsystem.Error.Printf("LogIn first")
		w.WriteHeader(http.StatusForbidden)
		view.WriteMessage(&w,view.ErrorMsg{"Login first"}, 1)
		return
	}
	md, errDb := user_model.UserInfoQuery(id.(int64))
	if errDb != nil {
		logsystem.Error.Printf("Database Error %s", errDb)
		w.WriteHeader(http.StatusInternalServerError)
		view.WriteMessage(&w,view.ErrorMsg{"Database Error"}, 2)
		return
	}
	inf := user.UserInfo{md.FamilyName.String, md.Name.String, md.SecondName.String,
		md.DateReceiving.String, md.IssuedBy.String, md.DivisionNumber.String,
		md.RegistrationAddres.String, md.MailingAddres.String, md.BirthDay.String,
		md.Sex.Bool, md.HomePhone.String, md.MobilePhone.String, md.CitizenShip.String,
	md.EMail}
	view.WriteMessage(&w,inf, 0)
}
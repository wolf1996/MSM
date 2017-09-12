package framework

import (
	"github.com/wolf1996/MSM/server/application/session_manager"
	"github.com/wolf1996/MSM/server/application/view"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/logsystem"
	"github.com/gorilla/sessions"
	"net/http"
	"context"
)

func SessionRequired(handlerFunc ContHandlerFunc) ContHandlerFunc{
	SessionRequiredHandler := func (appContext AppContext,w http.ResponseWriter, r *http.Request) {
		sess, err := session_manager.GetSession(r, "user_session")
		if err != nil {
			logsystem.Error.Printf("Get session error %s", err)
			sess, _ = session_manager.NewSession(r, "user_session")
			sess.Save(r, w)
		}
		cnt := r.Context()
		cnt = context.WithValue(cnt, "session", *sess)
		handlerFunc(appContext ,w,r.WithContext(cnt))
		return
	}
	return SessionRequiredHandler
}

func AuthRequired(handlerFunc ContHandlerFunc) ContHandlerFunc {
	AuthRequiredHandler := func (appContext AppContext, w http.ResponseWriter, r *http.Request) {
		cnt := r.Context()
		sess, ok := cnt.Value("session").(sessions.Session)
		if !ok {
			logsystem.Error.Printf("Middleware ERROR")
			w.WriteHeader(http.StatusInternalServerError)
			view.WriteMessage(&w, view.ErrorMsg{"Server Error"}, error_codes.NOT_LOGGED)
			return
		}
		id, ok := sess.Values["user"].(int64)
		if !ok {
			logsystem.Error.Printf("LogIn first")
			w.WriteHeader(http.StatusForbidden)
			view.WriteMessage(&w, view.ErrorMsg{"Login first"}, error_codes.NOT_LOGGED)
			return
		}
		cnt = context.WithValue(cnt, "id", id)
		handlerFunc(appContext, w,r.WithContext(cnt))
	}
	Handler := SessionRequired(AuthRequiredHandler)
	return Handler
}

func AppContextMiddleware(appContext AppContext, handlerFunc ContHandlerFunc) lowHandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		handlerFunc(appContext, w, r)
	}
}
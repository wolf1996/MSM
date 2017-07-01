package session_manager

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("somesecretkey"))

func GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

func NewSession(r *http.Request, name string) (*sessions.Session, error) {
	return store.New(r, name)
}

package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wolf1996/MSM/server/logsystem"
	"gopkg.in/mgo.v2"
)

var Database     sql.DB
var MongoSession *mgo.Session

func PostgresInit(userId, userPass, databaseURL string) error {
	logsystem.Info.Printf("Postgres login")
	cmd := fmt.Sprintf("postgres://%s:%s@%s", userId, userPass, databaseURL)
	Db, err := sql.Open("postgres", cmd)
	if err != nil {
		logsystem.Error.Printf("Database login error %s", err)
		return err
	}
	err = Db.Ping()
	if err != nil {
		return err
	}
	Database = *Db
	return nil
}

func MongoInit() error {
	logsystem.Info.Printf("Mongo login")
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		return  err
	}
	err = session.Ping()
	if err != nil {
		return  err
	}
	MongoSession = session
	return  nil
}
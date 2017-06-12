package models

import (
	"fmt"
	"github.com/wolf1996/MSM/server/logsystem"
	"database/sql"
	_"github.com/lib/pq"
)

var Database sql.DB

func Init(userId, userPass, databaseURL string) error{
	logsystem.Info.Printf("Database login start")
	cmd := fmt.Sprintf("postgres://%s:%s@%s", userId, userPass, databaseURL)
	Db, err := sql.Open("postgres", cmd)
	if err != nil{
		logsystem.Error.Printf("Database login error %s", err)
		return err
	}
	Database = *Db
	return nil
}


package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wolf1996/MSM/server/logsystem"
)

type MainDatabase  struct {
	Database *sql.DB
}

const DbSystemName string  = "database"


func GetDatabase(userId, userPass, databaseURL string)(database MainDatabase,err error){
	logsystem.Info.Printf("Database login start")
	cmd := fmt.Sprintf("postgres://%s:%s@%s", userId, userPass, databaseURL)
	Db, err := sql.Open("postgres", cmd)
	if err != nil {
		logsystem.Error.Printf("Database login error %s", err)
		return
	}
	database = MainDatabase{Db}
	return
}

package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"MSM/server/logsystem"
)

var Database sql.DB

func Init(userId, userPass, databaseURL string) error {
	logsystem.Info.Printf("Database login start")
	cmd := fmt.Sprintf("postgres://%s:%s@%s", userId, userPass, databaseURL)
	Db, err := sql.Open("postgres", cmd)
	if err != nil {
		logsystem.Error.Printf("Database login error %s", err)
		return err
	}
	Database = *Db
	return nil
}

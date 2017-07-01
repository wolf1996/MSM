package controller_model

import (
	"database/sql"
	_ "github.com/lib/pq"
)

//CREATE TABLE IF NOT EXISTS CONTROLLERS(
//id SERIAL PRIMARY KEY,
//name VARCHAR(256),
//user_id INT REFERENCES USERS(id),
//addres TEXT,
//activation_date DATE,
//status INT,
//mac MACADDR,
//deactivation_date DATE,
//controller_type INT
//);

type Table struct {
	Id               int
	Name             string
	UserId           int
	Adres            string
	ActivationDate   sql.NullString
	Status           int
	Mac              string
	DeactivationDate sql.NullString
	ControllerType   int
}

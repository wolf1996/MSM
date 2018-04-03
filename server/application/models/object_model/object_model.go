package object_model

import (
	_ "github.com/lib/pq"
)

//CREATE TABLE IF NOT EXISTS OBJECTS (
//id                SERIAL PRIMARY KEY,
//name              VARCHAR(256)              NOT NULL,
//user_id           INT REFERENCES USERS (id) NOT NULL,
//address           TEXT                      NOT NULL
//);

type Table struct {
	Id               int
	Name             string
	UserId           int
	Addres           string
}

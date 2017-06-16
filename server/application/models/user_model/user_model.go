package user_model

import (
	"time"
	"database/sql"
)

type Table struct{
	Id int64
	FamilyName sql.NullString
	Name sql.NullString
	SecondName sql.NullString
	DateReceiving time.Time
	IssuedBy sql.NullString
	DivisionNumber sql.NullString
	RegistrationAddres sql.NullString
	MailingAddres sql.NullString
	BirthDay time.Time
	Sex sql.NullBool
	HomePhone sql.NullString
	MobilePhone sql.NullString
	CitizenShip sql.NullString
	EMail string
	PassHash string
}
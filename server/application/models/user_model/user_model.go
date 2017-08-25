package user_model

import (
	"database/sql"
	"time"
)

type Table struct {
	Id                  int64
	FamilyName          sql.NullString
	Name                sql.NullString
	SecondName          sql.NullString
	DateReceiving       time.Time
	IssuedBy            sql.NullString
	DivisionNumber      sql.NullString
	RegistrationAddress sql.NullString
	MailingAddress      sql.NullString
	BirthDay            time.Time
	Sex                 sql.NullBool
	HomePhone           sql.NullString
	MobilePhone         sql.NullString
	CitizenShip         sql.NullString
	EMail               string
	PassHash            string
}

package user_model

import (
	"database/sql"
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type UserInfoModel struct {
	FamilyName          sql.NullString
	Name                sql.NullString
	SecondName          sql.NullString
	DateReceiving       sql.NullString
	IssuedBy            sql.NullString
	DivisionNumber      sql.NullString
	RegistrationAddress sql.NullString
	MailingAddress      sql.NullString
	BirthDay            sql.NullString
	Sex                 sql.NullBool
	HomePhone           sql.NullString
	MobilePhone         sql.NullString
	CitizenShip         sql.NullString
	EMail               string
}

/*
CREATE TABLE IF NOT EXISTS USERS(
  id SERIAL PRIMARY KEY,
  family_name VARCHAR(256),
  name VARCHAR(256),
  second_name VARCHAR(256) DEFAULT  NULL ,
  date_receiving DATE,
  issued_by TEXT,
  division_number VARCHAR(50),
  registration_address TEXT,
  mailing_address TEXT,
  home_phone VARCHAR(20),
  mobile_phone VARCHAR(20),
  citizenship VARCHAR(256),
  email VARCHAR(50),
  pass_hash VARCHAR(256)
);
*/

func UserInfoQuery(id int64) (UserInfoModel, models.ErrorModel) {
	qr, err := models.Database.Query(
		"SELECT family_name, name, second_name, date_receiving, issued_by, division_number, "+
			"registration_address ,mailing_address ,home_phone ,mobile_phone ,citizenship ,email "+
			"FROM USERS WHERE id = $1 ;", id)
	if err != nil {
		return UserInfoModel{}, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info UserInfoModel
	qr.Next()
	err = qr.Scan(&info.FamilyName, &info.Name, &info.SecondName, &info.DateReceiving, &info.IssuedBy,
		&info.DateReceiving, &info.RegistrationAddress, &info.MailingAddress, &info.HomePhone, &info.MobilePhone,
		&info.CitizenShip, &info.EMail)
	if err != nil {
		return UserInfoModel{}, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	return info, nil
}

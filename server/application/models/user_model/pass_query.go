package user_model

import (
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models"
	"strings"
)

type loginData struct {
	Id        int64
	UserEmail string
	PassHash  string
}

func LogInUser(user, pass string) (int64, models.ErrorModel) {
	qr, err := models.Database.Query("SELECT id, e_mail, pass_hash "+
		"FROM USERS WHERE e_mail = $1 ;", user)
	if err != nil {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var ldata loginData
	if !qr.Next() {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Invalid auth data %s", err), Code: error_codes.INVALID_AUTH_DATA}
	}
	err = qr.Scan(&ldata.Id, &ldata.UserEmail, &ldata.PassHash)
	if err != nil {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	if strings.Compare(ldata.PassHash, pass) != 0 {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Invalid auth data"), Code: error_codes.INVALID_AUTH_DATA}
	}
	return ldata.Id, nil
}

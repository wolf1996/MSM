package user_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
	"strings"
)

type loginData struct {
	Id int64
	UserEmail string
	PassHash string
}

func LogInUser(user, pass string)  (int64, models.ErrorModel){
	qr,err := models.Database.Query("SELECT id, e_mail, pass_hash " +
		"FROM USERS WHERE e_mail = $1 ;", user)
	if err != nil {
		return 0, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	var ldata loginData
	qr.Next()
	err = qr.Scan(&ldata.Id,&ldata.UserEmail,&ldata.PassHash)
	qr.Close()
	if err != nil {
		return 0, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	if strings.Compare(ldata.PassHash, pass) != 0 {
		return 0, models.ErrorModelImpl{Msg:fmt.Sprint("Invalid auth data"),Code:1}
	}
	return ldata.Id, nil
}
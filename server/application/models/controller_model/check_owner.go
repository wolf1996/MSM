package controller_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
)

func CheckIsOwner(userId, controllerId int64) (bool, models.ErrorModel){
	qr,err := models.Database.Query("SELECT id, e_mail, pass_hash " +
		"FROM CONTROLLER WHERE (id = $1) and (user_id = $2) ;", controllerId , userId)
	if err != nil {
		return false, models.ErrorModelImpl{Msg:fmt.Sprint("Database Error %s", err),Code:2}
	}
	defer qr.Close()
	if ! qr.Next(){
		return false, models.ErrorModelImpl{Msg:fmt.Sprint("Invalid owner %s", err),Code:1}
	}
	return true, nil
}

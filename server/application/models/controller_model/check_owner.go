package controller_model

import (
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

func CheckIsOwner(userId, controllerId int64) (bool, models.ErrorModel) {
	qr, err := models.Database.Query("SELECT id, email, pass_hash "+
		"FROM CONTROLLER WHERE (id = $1) and (user_id = $2) ;", controllerId, userId)
	if err != nil {
		return false, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	if !qr.Next() {
		return false, models.ErrorModelImpl{Msg: fmt.Sprint("Invalid owner %s", err), Code: error_codes.INVALID_OWNER}
	}
	return true, nil
}

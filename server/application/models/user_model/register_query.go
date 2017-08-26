package user_model

import (
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/view/user"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

func RegisterUser(form *user.RegisterForm) (int64, models.ErrorModel) {
	var qr, err = models.Database.Query("INSERT INTO USERS (" +
		"family_name, name, e_mail, pass_hash) VALUES ($1, $2, $3, $4) RETURNING id;",
		form.LastName, form.FirstName, form.Email, form.Password)
	if err != nil {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	if !qr.Next() {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Invalid auth data %s", err), Code: error_codes.INVALID_AUTH_DATA}
	}
	var id int64
	err = qr.Scan(&id)
	if err != nil {
		return 0, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	return id, nil
}
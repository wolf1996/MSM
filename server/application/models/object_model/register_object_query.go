package object_model
import (
	"github.com/wolf1996/MSM/server/application/models"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

func RegisterObjectQuery(userId int64, name,  addres string) models.ErrorModel {
	qr, err := models.Database.Query("INSERT INTO OBJECTS VALUES ( "+
		"default, $2, $1, $3 ) "+
		"RETURNING id", userId, name, addres)
	if err != nil {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Database Error ", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	if !qr.Next() {
		return models.ErrorModelImpl{Msg: fmt.Sprint("Error no controller with such id", err), Code: error_codes.DATABASE_INVALID_CONTROLLER}
	}
	return nil
}
package object_model

import (
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type ObjectModel Table

type ObjectModels []ObjectModel

func GetUserObjects(id int64) (ObjectModels, models.ErrorModel) {
	var infoSlice ObjectModels
	qr, err := models.Database.Query(
		"SELECT * "+
			"FROM OBJECTS WHERE user_id = $1 ;", id)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info ObjectModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.UserId, &info.Addres,)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}

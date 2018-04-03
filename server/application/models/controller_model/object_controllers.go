package controller_model

import (
	"fmt"
	"github.com/wolf1996/MSM/server/application/models"
	"github.com/wolf1996/MSM/server/application/error_codes"
)

type ControllerModel Table

type ControllerModels []ControllerModel

func GetObjectControllers(userId, id int64) (ControllerModels, models.ErrorModel) {
	var infoSlice ControllerModels
	qr, err := models.Database.Query(
		"SELECT CONTROLLERS.ID, CONTROLLERS.NAME, CONTROLLERS.OBJECT_ID, "+
			" CONTROLLERS.META, CONTROLLERS.ACTIVATION_DATE, CONTROLLERS.STATUS, "+
			" CONTROLLERS.MAC,  CONTROLLERS.DEACTIVATION_DATE,  CONTROLLERS.CONTROLLER_TYPE "+
			" FROM CONTROLLERS INNER JOIN OBJECTS ON CONTROLLERS.object_id = OBJECTS.id "+
				"WHERE object_id = $1 AND user_id = $2 ;", id, userId)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info ControllerModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.ObjectId, &info.Meta, &info.ActivationDate,
			&info.Status, &info.Mac, &info.DeactivationDate, &info.ControllerType)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}
/*
CREATE TABLE IF NOT EXISTS CONTROLLERS (
  id                SERIAL PRIMARY KEY,
  name              VARCHAR(256)              NOT NULL UNIQUE,
  object_id         INT REFERENCES OBJECTS (id) NOT NULL,
  meta              TEXT                      NOT NULL,
  activation_date   DATE DEFAULT NULL,
  status            INT  DEFAULT NULL,
  mac               MACADDR                   NOT NULL,
  deactivation_date DATE DEFAULT NULL,
  controller_type   INT  DEFAULT NULL
);
 */
func GetUserControllers(id int64) (ControllerModels, models.ErrorModel) {
	var infoSlice ControllerModels
	qr, err := models.Database.Query(
		"SELECT CONTROLLERS.ID, CONTROLLERS.NAME, CONTROLLERS.OBJECT_ID, "+
		" CONTROLLERS.META, CONTROLLERS.ACTIVATION_DATE, CONTROLLERS.STATUS, "+
		" CONTROLLERS.MAC,  CONTROLLERS.DEACTIVATION_DATE,  CONTROLLERS.CONTROLLER_TYPE"+
	" FROM CONTROLLERS INNER JOIN OBJECTS ON CONTROLLERS.object_id = OBJECTS.id WHERE user_id = $1 ;", id)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info ControllerModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.ObjectId, &info.Meta, &info.ActivationDate,
			&info.Status, &info.Mac, &info.DeactivationDate, &info.ControllerType)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}
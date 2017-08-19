package sensor_model

import (
	"database/sql"
	"fmt"
	"github.com/wolf1996/MSM/server/application/error_codes"
	"github.com/wolf1996/MSM/server/application/models"
)

type SensorModel struct {
	Id               int
	Name             string
	ControllerId     sql.NullInt64
	ActivationDate   sql.NullString
	Status           int
	DeactivationDate sql.NullString
	SensorType       int
	Company          string
}
type SensorModels []SensorModel

type SensorTaxedModel struct {
	Id      int64
	Name    string
	Type    int
	Status  int
	Tax     float64
	TaxName string
}
type SensorTaxedModels []SensorTaxedModel

func GetControlledSensors(controllerId, userId int64) (SensorModels, models.ErrorModel) {
	var infoSlice SensorModels
	qr, err := models.Database.Query(
		" SELECT "+
			" SENSOR.id, SENSOR.Name, SENSOR.controller_id, SENSOR.activation_date, SENSOR.status, SENSOR.deactivation_date, SENSOR.sensor_type, SENSOR.company"+
			" FROM SENSOR INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id WHERE controller_id = $1 AND user_id = $2;", controllerId, userId)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info SensorModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.ControllerId, &info.ActivationDate,
			&info.Status, &info.DeactivationDate, &info.SensorType, &info.Company)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}

func GetTaxedSensors(controllerId, userId int64) (SensorTaxedModels, models.ErrorModel) {
	var infoSlice SensorTaxedModels
	qr, err := models.Database.Query(
		" SELECT "+
			" SENSOR.id, SENSOR.Name, TAX.Tax, TAX.Name"+
			" FROM SENSOR INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id "+
			" JOIN TAX ON SENSOR.tax = TAX.id "+
			" WHERE controller_id = $1 AND user_id = $2;", controllerId, userId)
	if err != nil {
		return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()
	var info SensorTaxedModel
	for qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.Tax, &info.TaxName)
		if err != nil {
			return infoSlice, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
		infoSlice = append(infoSlice, info)
	}
	return infoSlice, nil
}

func GetTaxedSensor(sensorId, userId int64) (SensorTaxedModel, models.ErrorModel) {
	var info SensorTaxedModel
	qr, err := models.Database.Query(
		" SELECT "+
			" SENSOR.id, SENSOR.Name, SENSOR.sensor_type, SENSOR.status , TAX.Tax, TAX.Name"+
			" FROM SENSOR INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id "+
			" JOIN TAX ON SENSOR.tax = TAX.id "+
			" WHERE SENSOR.id = $1 AND user_id = $2;", sensorId, userId)
	if err != nil {
		return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
	}
	defer qr.Close()

	if qr.Next() {
		err = qr.Scan(&info.Id, &info.Name, &info.Type, &info.Status, &info.Tax, &info.TaxName)
		if err != nil {
			return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Error %s", err), Code: error_codes.DATABASE_ERROR}
		}
	} else {
		return info, models.ErrorModelImpl{Msg: fmt.Sprint("Database Invalid sensor"), Code: error_codes.DATABASE_INVALID_SENSOR}
	}
	return info, nil
}

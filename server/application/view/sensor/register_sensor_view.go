package sensor

type RegisterSensorForm struct {
	ControllerId int64 `json:"controller_id"`
	SensorId     int64 `json:"sensor_id"`
}

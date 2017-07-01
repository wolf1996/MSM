package sensor

type SensorInfo struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ControllerId     int64  `json:"controller_id"`
	ActivationDate   string `json:"activation_date"`
	Status           int    `json:"status"`
	DeactivationDate string `json:"deactivation_date"`
	SensorType       int    `json: "sensor_type"`
	Company          string `json: "company"`
}

type SensorsInfo []SensorInfo

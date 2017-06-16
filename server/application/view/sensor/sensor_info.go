package sensor


type SensorInfo struct {
	Id int
	Name string
	ControllerId int64
	ActivationDate string
	Status int
	DeactivationDate string
	SensorType int
	Company string
}

type SensorsInfo [] SensorInfo
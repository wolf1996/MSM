package controller

type ControllerInfo struct {
	Id               *int    `json:"id"`
	Name             *string `json:"name"`
	Object           *int    `json:"object_id"`
	Meta             *string `json:"meta"`
	ActivationDate   *string `json:"activation_date"`
	Status           *int    `json:"status"`
	Mac              *string `json:"mac"`
	DeactivationDate *string `json:"deactivation_date"`
	ControllerType   *int    `json:"controller_type"`
}

type ControllersInfo []ControllerInfo

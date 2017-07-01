package controller

type ControllerStats struct {
	CotrollerId int64   `json:"cotroller_id"`
	Month       float64 `json:"month"`
	PrevMonth   float64 `json:"prev_month"`
	PrevYear    float64 `json:"prev_year"`
}

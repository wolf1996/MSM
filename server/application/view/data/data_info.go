package data

import "time"

type DataInfo struct {
	SensorId int       `json:"sensor_id"`
	Date     time.Time `json:"date"`
	Value    int64     `json:"value"`
	Hash     string    `json:"hash"`
}

type DataInfoList []DataInfo

type DataInfoStats struct {
	CurrentMonth  float64 `json:"current_month"`
	PrevYearMonth float64 `json:"prev_year_month"`
	PrevYear      float64 `json:"prev_year_average"`
}

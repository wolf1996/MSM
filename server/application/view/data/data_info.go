package data

import (
	"time"
)

type DataInfo struct {
	SensorId *int       `json:"sensor_id"`
	Date     *time.Time `json:"date"`
	Value    *int64     `json:"value"`
	Hash     *string    `json:"hash"`
}

type DataInfoList []DataInfo

type DataInfoStats struct {
	CurrentMonth  *float64 `json:"current_month"`
	PrevYearMonth *float64 `json:"prev_year_month"`
	PrevYear      *float64 `json:"prev_year_average"`
}

type DataInfoStatsList []DataInfoStats

type SensorVidgetData struct {
	Type 		*int 		   `json:"type"`
	Name 		*string 	   `json:"name"`
	Status 		*int  		   `json:"status"`
	Accrual	 	*float32	   `json:"accrual"`
	Over 		*float32 	   `json:"over"`
	Result 		*float32 	   `json:"result"`
	Stats 		*DataInfoStats `json:"stats"`
}

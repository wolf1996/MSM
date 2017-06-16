package data

import "time"

type DataInfo struct {
	SensorId int `json:"sensor_id"`
	Date time.Time `json:"date"`
	Value int64 `json:"value"`
	Hash string `json:"hash"`
}

type DataInfoList []DataInfo

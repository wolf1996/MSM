package data_model

import "time"

type Table struct {
	SensorId int
	Date time.Time
	Value int64
	Hash string
}

/*
CREATE TABLE IF NOT EXISTS DATA(
  sensor_id INT REFERENCES SENSOR(id),
  date DATE,
  value BIGINT,
  hs UUID
);
 */
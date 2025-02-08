package services

import "github.com/influxdata/influxdb/pkg/snowflake"

var generator = snowflake.New(1)

func GenerateID() string {
	return generator.NextString()
}

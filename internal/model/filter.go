package model

import "time"

type Filter struct {
	StartDate time.Time
	EndDate   time.Time
	MinPrice  float64
	MaxPrice  float64
}

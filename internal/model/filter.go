package model

import "time"

type Filter struct {
	StartDate time.Time
	EndDate   time.Time
	MinPrice  int64
	MaxPrice  int64
}

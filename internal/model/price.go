package model

import "time"

type Price struct {
	ID         int64
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}

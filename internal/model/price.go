package model

import "time"

type Price struct {
	ID         int64
	Name       string
	Category   string
	Price      int64
	CreateDate time.Time
}

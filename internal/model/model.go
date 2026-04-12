package model

import "time"

type CattlePrice struct {
	Date     time.Time `json:"date"`
	Price    float64   `json:"price"`
	Unit     string    `json:"unit"`
	Category string    `json:"category"`
}

func (r CattlePrice) ParseDateToStr() string {
	return r.Date.Format("2006-01-02")
}

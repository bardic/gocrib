package model

import "time"

type StorePosition struct {
	Id        int `json:"-"`
	StoreId   string
	Lat       float32
	Long      float32
	Distance  float32
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

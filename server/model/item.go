package model

import "time"

type Item struct {
	Id                int `json:"-"`
	Name              string
	Cost              float32
	Weight            float32
	Unit              string
	Barcode           string
	StoreName         string
	StoreId           string
	StoreNeighborhood string
	Tags              []int
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

type Items struct {
	Items []Item
}

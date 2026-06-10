package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"type:numeric(10,2);not null"`
	Stock       int     `gorm:"default:0"`
	SellerID    uint
}

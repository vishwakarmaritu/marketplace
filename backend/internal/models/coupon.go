package models

import "gorm.io/gorm"

type Coupon struct {
	gorm.Model
	Code     string  `gorm:"unique;not null"`
	Discount float64 `gorm:"not null"`
	IsActive bool    `gorm:"default:true"`
}

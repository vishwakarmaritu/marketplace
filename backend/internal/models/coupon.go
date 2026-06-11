package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code       string    `gorm:"unique;not null" json:"code"`
	Discount   float64   `gorm:"not null" json:"discount"`
	ExpiryDate time.Time `json:"expiry_date"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
}

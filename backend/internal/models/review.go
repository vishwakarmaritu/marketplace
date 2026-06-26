package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	BuyerID   uint   `gorm:"not null" json:"buyer_id"`
	Rating    int    `gorm:"not null" json:"rating"`
	Comment   string `gorm:"type:text" json:"comment"`
}

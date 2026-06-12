package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"index;not null" json:"cart_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
}

package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	BuyerID     uint        `gorm:"not null" json:"buyer_id"`
	Status      string      `gorm:"type:varchar(20);default:'Placed';not null" json:"status"`
	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"index;not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
}

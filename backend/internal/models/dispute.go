package models

import "gorm.io/gorm"

type Dispute struct {
	gorm.Model
	OrderID       uint   `gorm:"index;not null" json:"order_id"`
	BuyerID       uint   `gorm:"not null" json:"buyer_id"`
	Reason        string `gorm:"type:text;not null" json:"reason"`
	Status        string `gorm:"type:varchar(20);default:'Open'" json:"status"`
	AdminResponse string `gorm:"type:text" json:"admin_response"`
}

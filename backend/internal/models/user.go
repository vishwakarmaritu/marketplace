package models

import (
	"errors"

	"gorm.io/gorm"
)

type Role string

const (
	RoleBuyer      Role = "buyer"
	RoleSeller     Role = "seller"
	RoleOperations Role = "operations"
	RoleAdmin      Role = "admin"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     Role   `gorm:"type:varchar(20);default:'buyer';not null" json:"role"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	switch u.Role {
	case RoleBuyer, RoleSeller, RoleOperations, RoleAdmin:
		return nil
	default:
		return errors.New("database error: invalid user role provided")
	}
}

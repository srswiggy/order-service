package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           int64 `gorm:"primaryKey"`
	UserUniqueID int64
	Location
}

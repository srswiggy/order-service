package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	ID           int64 `gorm:"primaryKey"`
	ItemUniqueID int64
	Name         string
	Price        float32
	OrderID      int64
}

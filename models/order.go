package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ID           int64 `gorm:"primaryKey"`
	RestaurantID int64
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID"`
	Status       string
	UserID       int64
	User         User `gorm:"foreignKey:UserID"`
	TotalPrice   float32
	Items        []MenuItem
}

package models

import "time"

type PurchasedItem struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	ItemName  string    `gorm:"not null"`
	Price     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Transaction struct {
	ID         uint      `gorm:"primaryKey"`
	FromUserID uint      `gorm:"not null"`
	ToUserID   uint      `gorm:"not null"`
	Amount     int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

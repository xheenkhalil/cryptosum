package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Wallet    string         `gorm:"uniqueIndex;not null" json:"wallet"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Token struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Symbol    string         `gorm:"not null" json:"symbol"`
	Address   string         `gorm:"uniqueIndex;not null" json:"address"`
	Chain     string         `gorm:"not null" json:"chain"`
	Decimals  int            `json:"decimals"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Portfolio struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Balance   float64        `json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

package model

import "time"

type Match struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	MatchedID uint      `json:"matched_id" gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Matched   bool      `gorm:"column=matched;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
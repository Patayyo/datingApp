package model

import "time"

type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user_id" gorm:"not null"`
	Bio       string `json:"bio"`
	Interests string `json:"interests"`
	Age       string `json:"age" gorm:"not null"`
	Gender    string `json:"gender" gorm:"not null"`
	Location  string `json:"location" gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
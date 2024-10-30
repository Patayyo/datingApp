package model

import "time"

// Match представляет модель для хранения информации о совпадениях(мэтчах) между пользователями в базе данных
type Match struct {
	ID        uint      `gorm:"primaryKey"` // Уникальный идентификатор мэтча(первичный ключ)
	UserID    uint      `json:"user_id" gorm:"not null"` // ID пользователя, инициировавшего мэтч
	MatchedID uint      `json:"matched_id" gorm:"not null"` // ID пользователя, с которым произошел мэтч
	User      User      `gorm:"foreignKey:UserID"` // Связь с моделью User через внешний ключ UserID
	Matched   bool      `gorm:"column=matched;not null"` // Флаг, указывающий, подтверждено ли совпадение
	CreatedAt time.Time `json:"created_at"` // Дата и время создания записи о матче
	UpdatedAt time.Time `json:"updated_at"` // Дата и время последнего обновления записи о мэтче
}
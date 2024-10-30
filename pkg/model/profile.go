package model

import "time"

// Profile представляет модель профиля пользователя в базе данных
// Содержит поля для хранения биографии, интересов и другой информации о пользователе
type Profile struct {
	ID        uint   `gorm:"primaryKey"` // Уникальный идентификатор профиля(первичный ключ)
	UserID    uint   `json:"user_id" gorm:"not null"` // ID пользователя, которому принадлежит профиль(обязательное поле)
	Bio       string `json:"bio"` // Биография пользователя
	Interests string `json:"interests"` // Интересы пользователя
	Age       string `json:"age" gorm:"not null"` // Возраст пользователя(обязательное поле)
	Gender    string `json:"gender" gorm:"not null"` // Пол пользователя(обязательное поле)
	Location  string `json:"location" gorm:"not null"`// Локация пользователя(обязательное поле)
	User      User   `gorm:"foreignKey:UserID"`// Связь с моделью User через внешний ключ UserID
	CreatedAt time.Time `json:"created_at"` // Дата и время создания профиля
	UpdatedAt time.Time `json:"updated_at"` // Дата и время последнего обновления профиля
}
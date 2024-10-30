package model

import (
	"time"
)

// User представляет модель пользователя в базе данных
// Содержит поля для хранения данных пользователя, таких как имя пользователя, email, пароль и связанная информация
type User struct {
	ID		 uint 	`gorm:"primaryKey"` // Уникальный идентификатор пользователя(первичный ключ)
	Username string `json:"username" gorm:"unique;not null"` // Имя пользователя(должно быть уникальным и не может быть пустым)
	Email    string `json:"email" gorm:"unique;not null"` // Электронная почта пользователя(уникальная и обязательная)
	Password string `json:"password" gorm:"not null"` // Пароль пользователя(обязательное поле)
	CreatedAt time.Time `json:"created_at"` // Дата и время создания записи о пользователе
	UpdatedAt time.Time `json:"updated_at"` // Дата и время последнего обновления записи о пользователе
	Profiles []Profile `json:"profiles" gorm:"foreignKey:UserID"` // Связь с профилями пользователя(внешний ключ UserID)
	Matches  []Match	`json:"matches" gorm:"foreignKey:UserID"` // Связь с записями о совпадениях(внешний ключ UserID) 
}
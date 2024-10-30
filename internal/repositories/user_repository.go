package repositories

import (
	"datingApp/pkg/model"

	"gorm.io/gorm"
)

// UserRepository представляет структуру, которая предоставляет методы для работы с пользователями и их профилями в базе данных.
type UserRepository struct {
	DB *gorm.DB // Ссылка на базу данных для выполнения операций
}

// GetUserByID получает пользователя по его уникальному идентификатору (userID).
func (r *UserRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	// Получаем первую запись из таблицы User с указанным userID
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, err // Возвращаем ошибку, если пользователь не найден
	}
	return &user, nil // Возвращаем найденного пользователя
}

// GetUserProfile получает профиль пользователя по его userID.
func (r *UserRepository) GetUserProfile(userID uint) (*model.Profile, error) {
	var profile model.Profile
	// Выполняем поиск в таблице Profile по полю user_id
	if err := r.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err // Возвращаем ошибку, если профиль не найден
	}
	return &profile, nil // Возвращаем найденный профиль пользователя 
}

// UpdateUserProfile обновляет информацию о пользователе и его профиле, используя userID.
func (r *UserRepository) UpdateUserProfile(userID uint, username, bio string) error {
	var user model.User
	// Находим пользователя по его userID
	if err := r.DB.First(&user, userID).Error; err != nil{
		return err // Возвращаем ошибку, если пользователь не найден
	}
	// Обновляем имя пользователя
	user.Username = username
	// Сохраняем изменения пользователя
	if err := r.DB.Save(&user).Error; err != nil{
		return err // Возвращаем ошибку, если обновление пользователя завершилось неудачно
	}

	var profile model.Profile
	// Находим профиль пользователя по userID
	if err := r.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return err // Возвращаем ошибку, если профиль не найден
	}
	// Обновляем поле "биография" в профиле
	profile.Bio = bio
	// Сохраняем изменения профиля
	return r.DB.Save(&profile).Error
}

// DeleteUser удаляет пользователя и его профиль по userID.
func (r *UserRepository) DeleteUser(userID uint) error {
	// Удаляем запись профиля, используя user_id
	if err := r.DB.Where("user_id = ?", userID).Delete(&model.Profile{}).Error; err != nil {
		return err // Возвращаем ошибку, если удаление профиля завершилось неудачно
	}
	// Удаляем пользователя с указанным userID
	return r.DB.Delete(&model.User{}, userID).Error 
}
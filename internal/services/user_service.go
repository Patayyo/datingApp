package services

import (
	"datingApp/pkg/model"
	"errors"

	"gorm.io/gorm"
)

// UserService представляет сервис для обработки операций с пользователями.
type UserService struct {
	DB *gorm.DB // Ссылка на базу данных для выполнения операций
}

// ProfileResponse представляет структуру ответа с информацией о профиле пользователя.
type ProfileResponse struct {
	Username string `json:"username"`
	Bio string `json:"bio"`
	Interests string `json:"interests"`
	Age string `json:"age"`
	Gender string `json:"gender"`
	Location string `json:"location"`
}

// UpdateProfileInput представляет входные данные для обновления профиля пользователя.
type UpdateProfileInput struct {
	Username *string `json:"username"`
	Bio *string `json:"bio"`
	Interests *string `json:"interests"`
	Age *string `json:"age"`
	Gender *string `json:"gender"`
	Location *string `json:"location"`
}

// GetUserProfile получает информацию о профиле пользователя по его ID.
func (s *UserService) GetUserProfile(userID int) (*ProfileResponse, error) {
	var profile model.Profile

	// Загружаем профиль пользователя с его связанной сущностью User
	if err := s.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found") // Возвращаем ошибку, если профиль не найден
		}
		return nil, err // Возвращаем другую ошибку, если произошла ошибка в запросе
	}

	// Формируем ответ с информацией о профиле
	response := &ProfileResponse{
		Username: profile.User.Username,
		Bio: profile.Bio,
		Interests: profile.Interests,
		Age: profile.Age,
		Gender: profile.Gender,
		Location: profile.Location,
	}
	return response, nil // Возвращаем ответ
}

// UpdateUserProfile обновляет информацию о пользователе и его профиле.
func (s *UserService) UpdateUserProfile(userID int, input UpdateProfileInput) error {
	var user model.User
	var profile model.Profile

	// Проверяем, существует ли пользователь
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found") // Возвращаем ошибку, если пользователь не найден
		}
		return err // Возвращаем другую ошибку, если произошла ошибка в запросе
	}

	// Если имя пользователя указано для обновления, обновляем его
	if input.Username != nil {
		user.Username = *input.Username
		if err := s.DB.Save(&user).Error; err != nil {
			return err // Возвращаем ошибку, если не удалось сохранить пользователя
		}
	}

	// Проверяем, существует ли профиль пользователя
	if err := s.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("profile not found") // Возвращаем ошибку, если профиль не найден
		}
		return err // Возвращаем другую ошибку, если произошла ошибка в запросе
	}

	// Обновляем поля профиля, если они указаны
	if input.Bio != nil {
		profile.Bio = *input.Bio
	}

	if input.Interests != nil {
		profile.Interests = *input.Interests
	}

	if input.Age != nil {
		profile.Age = *input.Age
	}

	if input.Gender != nil {
		profile.Gender = *input.Gender
	}

	if input.Location != nil {
		profile.Location = *input.Location
	}

	// Сохраняем обновленный профиль
	if err := s.DB.Save(&profile).Error; err != nil {
		return err // Возвращаем ошибку, если не удалось сохранить профиль
	}

	return nil // Возвращаем nil, если операция прошла успешно 
}

// DeleteUserProfile удаляет профиль и пользователя по его ID.
func (s *UserService) DeleteUserProfile(userID int) error {
	// Удаляем профиль пользователя
	if err := s.DB.Where("user_id = ?", userID).Delete(&model.Profile{}).Error; err != nil {
		return err // Возвращаем ошибку, если не удалось удалить профиль
	}

	// Удаляем пользователя
	if err := s.DB.Delete(&model.User{}, userID).Error; err != nil {
		return err // Возвращаем ошибку, если не удалось удалить пользователя
	}

	return nil // Возвращаем nil, если операция прошла успешно
}
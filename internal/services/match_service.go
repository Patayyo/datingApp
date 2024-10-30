package services

import (
	"datingApp/pkg/model"
	"errors"

	"gorm.io/gorm"
)

// MatchServiceInterface определяет интерфейс для взаимодействия с сервисом матчей.
type MatchServiceInterface interface {
    LikeUser(userID, targetID int) error // Метод для "лайка" пользователя
	GetUserMatches(userID int) ([]model.User, error) // Метод для получения списка мэтчей для пользователя
}
// MatchService представляет сервис для обработки операций с матчами. 
type MatchService struct {
	DB *gorm.DB // Ссылка на базу данных для выполнения операций
}

// LikeUser позволяет пользователю "лайкнуть" другого пользователя.
func (s *MatchService) LikeUser(userID int, targetUserID int) error {
	// Проверяем, не пытается ли пользователь "лайкнуть" себя
	if userID == targetUserID {
		return errors.New("you cannot like yourself")
	}

	// Проверяем, существует ли уже совпадение между пользователями
	var existingMatch model.Match
	if err := s.DB.Where("user_id = ? AND matched_id = ?", uint(userID), uint(targetUserID)).First(&existingMatch).Error; err == nil {
		return errors.New("you have already liked this user")
	}

	// Создаем новое совпадение
	newMatch := model.Match{
		UserID: uint(userID),
		MatchedID: uint(targetUserID),
	}

	// Сохраняем новое совпадение в базе данных
	if err := s.DB.Create(&newMatch).Error; err != nil {
		return err // Возвращаем ошибку, если не удалось сохранить совпадение
	}

	// Проверяем, существует ли обратное совпадение (лайк от целевого пользователя к текущему пользователю)
	var reciprocalMatch model.Match
	if err := s.DB.Where("user_id = ? AND matched_id = ?", uint(targetUserID), uint(userID)).First(&reciprocalMatch).Error; err == nil {
		// Если обратное совпадение существует, устанавливаем флаг "Matched" в true для обоих совпадений
		newMatch.Matched = true
		reciprocalMatch.Matched = true

		// Сохраняем обновленные совпадения в базе данных
		if err := s.DB.Save(&newMatch).Error; err != nil {
			return err
		}
		if err := s.DB.Save(&reciprocalMatch).Error; err != nil {
			return err
		}
	}

	return nil // Возвращаем nil, если операция прошла успешно
}

// GetUserMatches получает список пользователей, с которыми пользователь имеет совпадения.
func (s *MatchService) GetUserMatches(userID int) ([]model.User, error) {
	var matches []model.Match
	// Формируем SQL-запрос для получения всех совпадений текущего пользователя
	query := "SELECT * FROM matches WHERE user_id = ? AND matched = true"
	if err := s.DB.Raw(query, userID).Scan(&matches).Error; err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}

	// Извлекаем ID пользователей, с которыми есть совпадения
	var matchedUserIDs []uint
	for _, match := range matches {
		matchedUserIDs = append(matchedUserIDs, match.MatchedID)
	}

	// Получаем информацию о пользователях, с которыми есть совпадения
	var matchedUsers []model.User
	if err := s.DB.Where("id IN ?", matchedUserIDs).Find(&matchedUsers).Error; err != nil {
		return nil, err // Возвращаем ошибку, если не удалось найти пользователей
	}

	return matchedUsers, nil // Возвращаем список пользователей, с которыми есть совпадения
}
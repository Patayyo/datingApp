package repositories

import (
	"datingApp/pkg/model"

	"gorm.io/gorm"
)

// MatchRepository реализует методы для работы с сущностями Match в базе данных.
type MatchRepository struct {
	DB *gorm.DB // Ссылка на объект базы данных
}

// SaveMatch сохраняет новый матч в базе данных.
func (r *MatchRepository) SaveMatch(match *model.Match) error {
	// Создаем запись в таблице Match, возвращаем ошибку, если сохранение не удалось
	return r.DB.Create(match).Error
}

// GetUserMatches получает список матчей для конкретного пользователя по его userID.
func (r *MatchRepository) GetUserMatches(userID uint) ([]model.Match, error) {
	var matches []model.Match
	// Находим все матчи, в которых указан userID как user_id или matched_id
	err := r.DB.Where("user_id = ?", userID).Or("matched_id = ?", userID).Find(&matches).Error
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}
	return matches, nil // Возвращаем список матчей, если запрос выполнен успешно
}

// CheckIfMatched проверяет, существует ли матч между двумя пользователями.
func (r *MatchRepository) CheckIfMatched(userID, targetID uint) (bool, error) {
	var match model.Match
	// Проверяем наличие записи матча между userID и targetID в любом направлении
	err := r.DB.Where("(user_id = ? AND matched_id = ?) OR (user_id = ? AND matched_id = ?)", userID, targetID, targetID, userID).First(&match).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil // Если матч не найден, возвращаем false и nil ошибку
	}

	// Возвращаем true, если матч существует, и любую ошибку, если запрос завершился с ошибкой
	return err == nil, err 
}
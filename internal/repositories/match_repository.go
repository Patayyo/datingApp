package repositories

import (
	"datingApp/pkg/model"

	"gorm.io/gorm"
)

type MatchRepository struct {
	DB *gorm.DB
}

func (r *MatchRepository) SaveMatch(match *model.Match) error {
	return r.DB.Create(match).Error
}

func (r *MatchRepository) GetUserMatches(userID uint) ([]model.Match, error) {
	var matches []model.Match
	err := r.DB.Where("user_id = ?", userID).Or("matched_id = ?", userID).Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *MatchRepository) CheckIfMatched(userID, targetID uint) (bool, error) {
	var match model.Match
	err := r.DB.Where("(user_id = ? AND matched_id = ?) OR (user_id = ? AND matched_id = ?)", userID, targetID, targetID, userID).First(&match).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return err == nil, err
}
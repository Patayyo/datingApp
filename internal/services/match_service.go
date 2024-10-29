package services

import (
	"datingApp/pkg/model"
	"errors"

	"gorm.io/gorm"
)

type MatchServiceInterface interface {
    LikeUser(userID, targetID int) error
	GetUserMatches(userID int) ([]model.User, error)
}

type MatchService struct {
	DB *gorm.DB
}

func (s *MatchService) LikeUser(userID int, targetUserID int) error {
	if userID == targetUserID {
		return errors.New("you cannot like yourself")
	}

	var existingMatch model.Match
	if err := s.DB.Where("user_id = ? AND matched_id = ?", uint(userID), uint(targetUserID)).First(&existingMatch).Error; err == nil {
		return errors.New("you have already liked this user")
	}

	newMatch := model.Match{
		UserID: uint(userID),
		MatchedID: uint(targetUserID),
	}

	if err := s.DB.Create(&newMatch).Error; err != nil {
		return err
	}

	var reciprocalMatch model.Match
	if err := s.DB.Where("user_id = ? AND matched_id = ?", uint(targetUserID), uint(userID)).First(&reciprocalMatch).Error; err == nil {
		newMatch.Matched = true
		reciprocalMatch.Matched = true

		if err := s.DB.Save(&newMatch).Error; err != nil {
			return err
		}
		if err := s.DB.Save(&reciprocalMatch).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *MatchService) GetUserMatches(userID int) ([]model.User, error) {
	var matches []model.Match
	query := "SELECT * FROM matches WHERE user_id = ? AND matched = true"
	if err := s.DB.Raw(query, userID).Scan(&matches).Error; err != nil {
		return nil, err
	}

	var matchedUserIDs []uint
	for _, match := range matches {
		matchedUserIDs = append(matchedUserIDs, match.MatchedID)
	}

	var matchedUsers []model.User
	if err := s.DB.Where("id IN ?", matchedUserIDs).Find(&matchedUsers).Error; err != nil {
		return nil, err
	}

	return matchedUsers, nil
}
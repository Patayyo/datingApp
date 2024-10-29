package services

import (
	"datingApp/pkg/model"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

type ProfileResponse struct {
	Username string `json:"username"`
	Bio string `json:"bio"`
	Interests string `json:"interests"`
	Age string `json:"age"`
	Gender string `json:"gender"`
	Location string `json:"location"`
}

type UpdateProfileInput struct {
	Username *string `json:"username"`
	Bio *string `json:"bio"`
	Interests *string `json:"interests"`
	Age *string `json:"age"`
	Gender *string `json:"gender"`
	Location *string `json:"location"`
}

func (s *UserService) GetUserProfile(userID int) (*ProfileResponse, error) {
	var profile model.Profile

	if err := s.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	response := &ProfileResponse{
		Username: profile.User.Username,
		Bio: profile.Bio,
		Interests: profile.Interests,
		Age: profile.Age,
		Gender: profile.Gender,
		Location: profile.Location,
	}
	return response, nil
}

func (s *UserService) UpdateUserProfile(userID int, input UpdateProfileInput) error {
	var user model.User
	var profile model.Profile

	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if input.Username != nil {
		user.Username = *input.Username
		if err := s.DB.Save(&user).Error; err != nil {
			return err 
		}
	}

	if err := s.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("profile not found")
		}
		return err
	}

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

	if err := s.DB.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUserProfile(userID int) error {
	if err := s.DB.Where("user_id = ?", userID).Delete(&model.Profile{}).Error; err != nil {
		return err
	}

	if err := s.DB.Delete(&model.User{}, userID).Error; err != nil {
		return err
	}

	return nil
}
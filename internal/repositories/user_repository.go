package repositories

import (
	"datingApp/pkg/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserProfile(userID uint) (*model.Profile, error) {
	var profile model.Profile
	if err := r.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) UpdateUserProfile(userID uint, username, bio string) error {
	var user model.User
	if err := r.DB.First(&user, userID).Error; err != nil{
		return err
	}
	user.Username = username
	if err := r.DB.Save(&user).Error; err != nil{
		return err 
	}

	var profile model.Profile
	if err := r.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return err
	}
	profile.Bio = bio
	return r.DB.Save(&profile).Error
}

func (r *UserRepository) DeleteUser(userID uint) error {
	if err := r.DB.Where("user_id = ?", userID).Delete(&model.Profile{}).Error; err != nil {
		return err
	}
	return r.DB.Delete(&model.User{}, userID).Error
}
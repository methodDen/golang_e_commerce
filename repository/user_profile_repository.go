package repository

import (
	"E-Commerce/models"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	DB *gorm.DB
}

type UserProfileRepository interface {
	GetUserProfile(int) (models.UserProfile, error)
	UpdateUserProfile(models.UserProfile) (models.UserProfile, error)
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return userProfileRepository{
		DB: db,
	}
}

// GetUserProfile get profile by user_id (foreign key)
func (u userProfileRepository) GetUserProfile(userID int) (userProfile models.UserProfile, err error) {
	return userProfile, u.DB.Where("user_id = ?", userID).First(&userProfile).Error
}

func (u userProfileRepository) UpdateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	if err := u.DB.Where("user_id = ?", userProfile.UserID).First(&userProfile).Error; err != nil {
		return userProfile, err
	}
	return userProfile, u.DB.Model(&userProfile).Updates(&userProfile).Error
}

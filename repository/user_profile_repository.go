package repository

import (
	"E-Commerce/models"
	"E-Commerce/serializers/request"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	DB *gorm.DB
}

type UserProfileRepository interface {
	GetUserProfile(int) (models.UserProfile, error)
	UpdateUserProfile(models.UserProfile, request.UpdateUserProfileInput) (models.UserProfile, error)
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

func (u userProfileRepository) UpdateUserProfile(
	userProfile models.UserProfile,
	input request.UpdateUserProfileInput,
) (models.UserProfile, error) {
	userProfile.FirstName = input.FirstName
	userProfile.LastName = input.LastName
	userProfile.Country = input.Country
	userProfile.City = input.City
	userProfile.Address = input.Address
	return userProfile, u.DB.Save(&userProfile).Error
}

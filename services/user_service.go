package services

import (
	"errors"
	"strings"

	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/dto"
	"github.com/maycolacerda/ticketfair/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	var existing models.User

	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already in use")
	}
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, errors.New("username already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := models.User{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Username: strings.TrimSpace(req.Username),
		Password: string(hash),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &dto.UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func GetAllUsers(page, limit int) (*dto.PaginatedUsersResponse, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, errors.New("failed to count users")
	}

	if err := database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users")
	}

	data := make([]dto.UserResponse, len(users))
	for i, u := range users {
		data[i] = dto.UserResponse{
			UserID:    u.UserID,
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}
	}

	return &dto.PaginatedUsersResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func GetUserByID(userID string) (*dto.UserResponse, error) {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

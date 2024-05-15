package services

import "JobquestApi/models"

type UserService interface {
	CreateUser(*models.User) error
	GetUser(*string) (*models.User, error)
	GetAll()(users []*models.User, err error)
	UpdateUser(email string, user *models.User)(*models.User, error)
	DeleteUser(*string) error
}
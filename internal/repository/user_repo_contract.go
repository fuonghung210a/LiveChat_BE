package repository

import "go_starter/internal/model"

type IUserRepository interface {
	Create(user *model.User) error
	FindAll() ([]*model.User, error)
	FindById(userId int64) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	UpdateById(userId int64, user *model.User) error
	DeleteById(userId int64) error
	Paginate(page int32, pageSize int32) ([]*model.User, error)
}

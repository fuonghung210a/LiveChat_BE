package service

import (
	"go_starter/internal/model"
	"go_starter/internal/repository"
	"go_starter/internal/util"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email, password string) (*model.User, error) {
	// Hash the password before storing
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{Name: name, Email: email, Password: hashedPassword}
	err = s.repo.Create(user)
	return user, err
}

func (s *UserService) ListUsers() ([]*model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserById(userId int64) (*model.User, error) {
	return s.repo.FindById(userId)
}

func (s *UserService) UpdateUser(userId int64, name, email, password string) (*model.User, error) {
	// Hash the password before updating
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{Name: name, Email: email, Password: hashedPassword}
	err = s.repo.UpdateById(userId, user)
	if err != nil {
		return nil, err
	}
	return s.repo.FindById(userId)
}

func (s *UserService) DeleteUser(userId int64) error {
	return s.repo.DeleteById(userId)
}

func (s *UserService) PaginateUsers(page int32, pageSize int32) ([]*model.User, error) {
	return s.repo.Paginate(page, pageSize)
}

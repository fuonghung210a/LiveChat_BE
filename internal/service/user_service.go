package service

import (
	"go_starter/internal/model"
	"go_starter/internal/repository"
	"go_starter/internal/util"
)

type UserService struct {
	repo         *repository.UserRepository
	emailService *EmailService
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo:         repo,
		emailService: NewEmailService(),
	}
}

func (s *UserService) CreateUser(name, email, password string) (*model.User, error) {
	// Hash the password before storing
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{Name: name, Email: email, Password: hashedPassword}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Send welcome email (async to not block the response)
	go func() {
		if err := s.emailService.SendWelcomeEmail(user.Email, user.Name); err != nil {
			// Log the error but don't fail the user creation
			println("Failed to send welcome email: ", err.Error())
		}
	}()

	return user, nil
}

func (s *UserService) ListUsers() ([]*model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserById(userId int64) (*model.User, error) {
	return s.repo.FindById(userId)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
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

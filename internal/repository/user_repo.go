package repository

import (
	"go_starter/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) FindById(userId int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, userId).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateById(userId int64, user *model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", userId).Updates(user).Error
}

func (r *UserRepository) DeleteById(userId int64) error {
	return r.db.Delete(&model.User{}, userId).Error
}

func (r *UserRepository) Paginate(page int32, pageSize int32) ([]*model.User, error) {
	var users []*model.User
	offset := (page - 1) * pageSize
	err := r.db.Offset(int(offset)).Limit(int(pageSize)).Find(&users).Error
	return users, err
}

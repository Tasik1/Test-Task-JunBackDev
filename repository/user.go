package repository

import (
	"TestBackDev/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(model.User) (model.User, error)
	GetByLogin(string) (model.User, error)
	StoreRefreshToken(model.Token) error
	GetTokenByUserID(uint) (model.Token, error)
	UpdateRefreshToken(model.Token) error
}

type userRepository struct {
	connection *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		connection: DB(),
	}
}

func (db *userRepository) CreateUser(user model.User) (model.User, error) {
	return user, db.connection.Create(&user).Error
}

func (db *userRepository) GetByLogin(login string) (user model.User, err error) {
	return user, db.connection.First(&user, "login=?", login).Error
}

func (db *userRepository) GetTokenByUserID(id uint) (token model.Token, err error) {
	return token, db.connection.First(&token, "guid=?", id).Error
}

func (db *userRepository) StoreRefreshToken(token model.Token) (err error) {
	return db.connection.Create(&token).Error
}

func (db *userRepository) UpdateRefreshToken(token model.Token) error {
	if err := db.connection.First(&model.Token{}, token.GUID).Error; err != nil {
		return err
	}
	return db.connection.Model(&token).Updates(&token).Error
}

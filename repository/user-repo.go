package repository

import (
	"context"
	"mods/entity"

	"gorm.io/gorm"
)

type userConnection struct {
	connection *gorm.DB
}

type UserRepository interface {
	// functional
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	CheckNoId(ctx context.Context, NoId string) (entity.User, bool, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (db *userConnection) GetAllUser(ctx context.Context) ([]entity.User, error) {
	var listUser []entity.User

	tx := db.connection.Find(&listUser)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return listUser, nil
}

func (db *userConnection) CheckNoId(ctx context.Context, NoId string) (entity.User, bool, error) {

	var user entity.User
	if err := db.connection.Find(&user).Where("noid = ?", NoId).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

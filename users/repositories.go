package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(context *gin.Context, condition interface{}) (UserModel, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (repo userRepositoryImpl) FindUser(ctx *gin.Context, condition interface{}) (UserModel, error) {
	var model UserModel
	err := repo.db.WithContext(ctx).Where(condition).First(&model).Error
	return model, err
}

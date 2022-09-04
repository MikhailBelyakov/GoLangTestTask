package users

import (
	"github.com/gin-gonic/gin"
	"testProject/common"
)

func (c UserModel) TableName() string {
	return "users"
}

type UserModel struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"column:username"`
}

func AutoMigrate() error {
	db := common.GetDB()
	err := db.AutoMigrate(&UserModel{})
	if err != nil {
		return err
	}
	return nil
}

func FindOneUser(context *gin.Context, condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

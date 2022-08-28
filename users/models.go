package users

import (
	"testGoProject/common"
)

type Tabler interface {
	TableName() string
}

func (c UserModel) TableName() string {
	return "users"
}

type UserModel struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"column:username"`
}

func AutoMigrate() {
	db := common.GetDB()
	err := db.AutoMigrate(&UserModel{})
	if err != nil {
		return
	}
}

func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

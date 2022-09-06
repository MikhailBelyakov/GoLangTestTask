package users

import (
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

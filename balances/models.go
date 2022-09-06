package balances

import (
	"gorm.io/gorm"
	"testProject/common"
	"testProject/users"
)

func (model *BalanceModel) TableName() string {
	return "balances"
}

type BalanceModel struct {
	gorm.Model
	User   users.UserModel
	UserID uint   `gorm:"unique_index"`
	Amount uint32 `json:"amount"`
}

func AutoMigrate() error {
	db := common.GetDB()
	err := db.AutoMigrate(&BalanceModel{})
	if err != nil {
		return err
	}
	return nil
}

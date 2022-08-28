package balances

import (
	"github.com/jinzhu/gorm"
	"testGoProject/common"
	"testGoProject/users"
)

type Tabler interface {
	TableName() string
}

func (model *BalanceModel) TableName() string {
	return "balances"
}

type BalanceModel struct {
	gorm.Model
	User   users.UserModel
	UserID uint   `gorm:"unique_index"`
	Amount uint32 `json:"amount"`
}

func AutoMigrate() {
	db := common.GetDB()
	err := db.AutoMigrate(&BalanceModel{})
	if err != nil {
		return
	}
}

func (model *BalanceModel) changeBalance(amount uint32) error {
	db := common.GetDB()
	err := db.Model(model).Update("amount", amount).Error
	return err
}

// 	Получаем баланс, если не найден - создаём.
func GetBalance(userModel users.UserModel) BalanceModel {
	var balanceModel BalanceModel

	db := common.GetDB()
	db.Select("id", "user_id", "amount").Where(&BalanceModel{
		UserID: userModel.ID,
	}).FirstOrCreate(&balanceModel)

	return balanceModel
}

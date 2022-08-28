package transactions

import (
	"github.com/jinzhu/gorm"
	"testGoProject/common"
	"testGoProject/users"
	"time"
)

type Tabler interface {
	TableName() string
}

func (model *TransactionModel) TableName() string {
	return "transactions"
}

type TransactionModel struct {
	gorm.Model
	UserID     uint
	ReceiverId uint `gorm:"default: null"`
	SenderId   uint `gorm:"default: null"`
	Amount     uint32
	Date       time.Time
	Operation  string
}

func AutoMigrate() {
	db := common.GetDB()
	db.AutoMigrate(&TransactionModel{})
}

// 	Получаем транзакции пользователя
func FindTransactionsByUser(userModel users.UserModel) []TransactionModel {
	var transactionModels []TransactionModel

	db := common.GetDB()
	db.Where(&TransactionModel{
		UserID: userModel.ID,
	}).Find(&transactionModels)

	return transactionModels
}

func CreateTransaction(model TransactionModel) error {
	db := common.GetDB()
	err := db.Save(&model).Error
	return err
}

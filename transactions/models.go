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
	TypeID     int
	ReceiverId uint `gorm:"default: null"`
	SenderId   uint `gorm:"default: null"`
	Amount     uint32
	Date       time.Time
}

func AutoMigrate() {
	db := common.GetDB()
	err := db.AutoMigrate(&TransactionModel{})
	if err != nil {
		return
	}
}

// 	Получаем транзакции пользователя
func FindTransactionsByUser(userModel users.UserModel, condition GetTransactionParamStruct) []TransactionModel {
	var transactionModels []TransactionModel

	db := common.GetDB()
	db.Where(&TransactionModel{
		UserID: userModel.ID,
	}).Limit(condition.Limit).Order(condition.Order + " " + condition.Sort).Offset(condition.Offset).Find(&transactionModels)

	return transactionModels
}

func CreateTransaction(model TransactionModel) error {
	db := common.GetDB()
	err := db.Save(&model).Error
	return err
}

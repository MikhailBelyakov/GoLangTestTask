package transactions

import (
	"gorm.io/gorm"
	"testProject/common"
	"time"
)

func (model *TransactionModel) TableName() string {
	return "transactions"
}

type TransactionModel struct {
	gorm.Model
	UserID     uint
	TypeID     int
	ReceiverID uint `gorm:"default: null"`
	SenderID   uint `gorm:"default: null"`
	Amount     uint32
	Date       time.Time
}

func AutoMigrate() error {
	db := common.GetDB()
	err := db.AutoMigrate(&TransactionModel{})
	if err != nil {
		return err
	}
	return nil
}

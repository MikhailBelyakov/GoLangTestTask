package transactions

import (
	"github.com/gin-gonic/gin"
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

type TransactionRepository interface {
	FindByUserID(context *gin.Context, userID int, condition GetTransactionParamStruct) ([]TransactionModel, error)
	CreateTransaction(context *gin.Context, model *TransactionModel) error
}

func NewTransactionRepository(db *gorm.DB) (TransactionRepository, error) {
	return &transactionRepositoryImpl{
		db: db,
	}, nil
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func (repo *transactionRepositoryImpl) CreateTransaction(ctx *gin.Context, model *TransactionModel) error {
	err := repo.db.WithContext(ctx).Save(&model).Error
	return err
}

func (repo *transactionRepositoryImpl) FindByUserID(ctx *gin.Context, userID int, condition GetTransactionParamStruct) ([]TransactionModel, error) {
	var transaction []TransactionModel

	repo.db.WithContext(ctx).Where(&TransactionModel{
		UserID: uint(userID),
	}).Limit(condition.Limit).Order(condition.Order + " " + condition.Sort).Offset(condition.Offset).Find(&transaction)

	return transaction, nil
}

package transactions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByUserID(context *gin.Context, userID int, condition GetTransactionParamStruct) ([]TransactionModel, error)
	CreateTransaction(context *gin.Context, model *TransactionModel) error
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
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

package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testProject/common"
)

type TransactionService interface {
	GetTransactionsByUser(ctx *gin.Context, userID int, paramStruct GetTransactionParamStruct) ([]TransactionModel, common.HttpError)
}

func NewTransactionService(repo TransactionRepository) TransactionService {
	return &transactionsServiceImpl{
		repo: repo,
	}
}

type transactionsServiceImpl struct {
	repo TransactionRepository
}

func (service *transactionsServiceImpl) GetTransactionsByUser(ctx *gin.Context, userID int, paramStruct GetTransactionParamStruct) ([]TransactionModel, common.HttpError) {
	transactions, err := service.repo.FindByUserID(ctx, userID, paramStruct)
	if err != nil {
		return nil, common.NewHttpError(http.StatusNotFound, err)
	}

	return transactions, nil
}

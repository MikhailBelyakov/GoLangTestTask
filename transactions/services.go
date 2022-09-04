package transactions

import "github.com/gin-gonic/gin"

type TransactionService interface {
	GetTransactionsByUser(ctx *gin.Context, userID int, paramStruct GetTransactionParamStruct) ([]TransactionModel, error)
}

func NewTransactionService(repo TransactionRepository) TransactionService {
	return &transactionsServiceImpl{
		repo: repo,
	}
}

type transactionsServiceImpl struct {
	repo TransactionRepository
}

func (s *transactionsServiceImpl) GetTransactionsByUser(ctx *gin.Context, userID int, paramStruct GetTransactionParamStruct) ([]TransactionModel, error) {
	transactions, err := s.repo.FindByUserID(ctx, userID, paramStruct)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

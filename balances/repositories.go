package balances

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BalanceRepository interface {
	UpdateBalance(context *gin.Context, balanceModel *BalanceModel) error
	GetBalance(context *gin.Context, userID uint) BalanceModel
}

type balanceRepositoryImpl struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepositoryImpl{
		db: db,
	}
}

func (repo balanceRepositoryImpl) UpdateBalance(ctx *gin.Context, balanceModel *BalanceModel) error {
	err := repo.db.WithContext(ctx).Model(balanceModel).Update("amount", balanceModel.Amount).Error
	return err
}

func (repo balanceRepositoryImpl) GetBalance(ctx *gin.Context, userID uint) BalanceModel {
	var balanceModel BalanceModel

	repo.db.WithContext(ctx).Select("id", "user_id", "amount").Where(&BalanceModel{
		UserID: userID,
	}).FirstOrCreate(&balanceModel)

	return balanceModel
}

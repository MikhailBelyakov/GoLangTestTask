package balances

import (
	"github.com/gin-gonic/gin"
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

type BalanceRepository interface {
	UpdateBalance(context *gin.Context, balanceModel *BalanceModel) error
	GetBalance(context *gin.Context, userID uint) BalanceModel
}

func NewBalanceRepository(db *gorm.DB) (BalanceRepository, error) {
	return &balanceRepositoryImpl{
		db: db,
	}, nil
}

type balanceRepositoryImpl struct {
	db *gorm.DB
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

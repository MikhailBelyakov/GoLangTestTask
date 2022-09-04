package balances

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type ExchangeParamStruct struct {
	senderID   uint
	receiverID uint
	amount     float64
}

func (input *ExchangeParamStruct) BindExchangeParams(context *gin.Context) error {
	var err error

	senderIDString := context.Param("userID")
	receiverIDString := context.Param("receiverID")
	amountParam := context.PostForm("amount")

	senderID, err := strconv.ParseUint(senderIDString, 10, 32)
	receiverID, err := strconv.ParseUint(receiverIDString, 10, 32)
	amount, err := strconv.ParseFloat(amountParam, 64)

	if err != nil {
		return err
	}

	input.senderID = uint(senderID)
	input.receiverID = uint(receiverID)
	input.amount = amount

	return nil
}

type ChangeParamStruct struct {
	userID uint
	amount float64
}

func (input *ChangeParamStruct) BindChangeParams(context *gin.Context) error {
	var err error

	userIDString := context.Param("userID")
	amountParam := context.PostForm("amount")

	userID, err := strconv.ParseUint(userIDString, 10, 32)
	amount, err := strconv.ParseFloat(amountParam, 64)

	if err != nil {
		return err
	}

	input.userID = uint(userID)
	input.amount = amount

	return nil
}

type GetBalanceParamStruct struct {
	userID   uint
	currency string
}

func (input *GetBalanceParamStruct) BindGetBalanceParams(context *gin.Context) error {
	var err error

	userIDString := context.Param("userID")
	currency := context.Query("currency")

	userID, err := strconv.ParseInt(userIDString, 10, 32)

	if err != nil {
		return err
	}

	input.userID = uint(userID)
	input.currency = currency

	return nil
}

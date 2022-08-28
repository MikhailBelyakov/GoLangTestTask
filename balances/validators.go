package balances

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type ExchangeParamStruct struct {
	senderId   uint
	receiverId uint
	amount     float64
}

func (input *ExchangeParamStruct) BindExchangeParams(context *gin.Context) error {
	var err error

	senderIdString := context.Param("userId")
	receiverIdString := context.Param("receiverId")
	amountParam := context.PostForm("amount")

	senderId, err := strconv.ParseUint(senderIdString, 10, 32)
	receiverId, err := strconv.ParseUint(receiverIdString, 10, 32)
	amount, err := strconv.ParseFloat(amountParam, 64)

	if err != nil {
		return err
	}

	input.senderId = uint(senderId)
	input.receiverId = uint(receiverId)
	input.amount = amount

	return nil
}

type ChangeParamStruct struct {
	userId uint
	amount float64
}

func (input *ChangeParamStruct) BindChangeParams(context *gin.Context) error {
	var err error

	userIdString := context.Param("userId")
	amountParam := context.PostForm("amount")

	userId, err := strconv.ParseUint(userIdString, 10, 32)
	amount, err := strconv.ParseFloat(amountParam, 64)

	if err != nil {
		return err
	}

	input.userId = uint(userId)
	input.amount = amount

	return nil
}

type GetBalanceParamStruct struct {
	userId   uint
	currency string
}

func (input *GetBalanceParamStruct) BindGetBalanceParams(context *gin.Context) error {
	var err error

	userIdString := context.Param("userId")
	currency := context.Query("currency")

	userId, err := strconv.ParseInt(userIdString, 10, 32)

	if err != nil {
		return err
	}

	input.userId = uint(userId)
	input.currency = currency

	return nil
}

package transactions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"testGoProject/common"
	"testGoProject/users"
)

type GetTransactionParamStruct struct {
	userId uint
	common.SortStruct
	common.PaginateStruct
}

func orderFieldAllow() map[string]bool {
	return map[string]bool{
		"amount": true,
		"date":   true,
	}
}

func (input *GetTransactionParamStruct) BindParam(context *gin.Context) error {
	var err error

	input.PaginateStruct.BindParam(context)
	input.SortStruct.BindParam(context)

	userIdString := context.Param("userId")

	userId, err := strconv.ParseUint(userIdString, 10, 32)

	if err != nil {
		return errors.New(users.UserNotAvailableParamMessage)
	}

	if _, inMap := orderFieldAllow()[input.Order]; !inMap {
		return errors.New(common.SortFieldErrorMessage)
	}

	input.userId = uint(userId)

	return nil
}

package transactions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testProject/common"
	"testProject/users"
)

type GetTransactionParamStruct struct {
	userID uint
	common.SortStruct
	common.PaginateStruct
}

func orderFieldAllow() map[string]bool {
	return map[string]bool{
		"amount": true,
		"date":   true,
	}
}

func (input *GetTransactionParamStruct) BindParam(context *gin.Context) common.HttpError {
	var err error

	input.PaginateStruct.BindParam(context)
	input.SortStruct.BindParam(context)

	userIDString := context.Param("userID")

	userID, err := strconv.ParseUint(userIDString, 10, 32)

	if err != nil {
		return common.NewHttpError(http.StatusBadRequest, errors.New(users.UserNotAvailableParamMessage))
	}

	if _, inMap := orderFieldAllow()[input.Order]; !inMap {
		return common.NewHttpError(http.StatusBadRequest, errors.New(common.SortFieldErrorMessage))
	}

	input.userID = uint(userID)

	return nil
}

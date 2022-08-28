package transactions

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type GetTransactionParamStruct struct {
	userId uint
}

func (input *GetTransactionParamStruct) BindGetTransactionParam(context *gin.Context) error {
	var err error

	userIdString := context.Param("userId")

	userId, err := strconv.ParseUint(userIdString, 10, 32)

	if err != nil {
		return err
	}

	input.userId = uint(userId)

	return nil
}

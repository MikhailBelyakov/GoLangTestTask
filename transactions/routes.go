package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserTransaction(router *gin.RouterGroup) {
	router.GET("/:userId", getUserTransaction)
}

func getUserTransaction(context *gin.Context) {
	var inputParamStruct GetTransactionParamStruct
	var transactions []TransactionModel

	err := inputParamStruct.BindGetTransactionParam(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	transactions, err = GetTransactionsByUser(inputParamStruct)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": userNotFoundMessage})
		return
	}

	serializer := TransactionsStruct{context, transactions}
	context.JSON(http.StatusOK, gin.H{"data": serializer.Response()})
}

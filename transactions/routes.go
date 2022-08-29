package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testGoProject/users"
)

func UserTransaction(router *gin.RouterGroup) {
	router.GET("/:userId", getUserTransaction)
}

func getUserTransaction(context *gin.Context) {
	var inputParamStruct GetTransactionParamStruct
	var transactions []TransactionModel

	err := inputParamStruct.BindParam(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	transactions, err = GetTransactionsByUser(inputParamStruct)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": users.UserNotFoundMessage})
		return
	}

	serializer := TransactionsStruct{context, transactions}
	context.JSON(http.StatusOK, gin.H{"data": serializer.Response(), "paginate": inputParamStruct.PaginateStruct, "order": inputParamStruct.SortStruct})
}

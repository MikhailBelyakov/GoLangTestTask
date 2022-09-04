package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testProject/users"
)

func UserTransaction(router *gin.RouterGroup, controller TransactionController) {
	router.GET("/:userID", controller.GetUserTransaction)
}

type TransactionController interface {
	GetUserTransaction(ctx *gin.Context)
}

func NewTransactionController(service TransactionService) TransactionController {
	return &transactionControllerImpl{
		service: service,
	}
}

type transactionControllerImpl struct {
	service TransactionService
}

func (controller *transactionControllerImpl) GetUserTransaction(ctx *gin.Context) {
	var inputParamStruct GetTransactionParamStruct
	var transactions []TransactionModel

	err := inputParamStruct.BindParam(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	transactions, err = controller.service.GetTransactionsByUser(ctx, int(inputParamStruct.userID), inputParamStruct)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": users.UserNotFoundMessage})
		return
	}

	serializer := TransactionsStruct{ctx, transactions}
	ctx.JSON(http.StatusOK, gin.H{"data": serializer.Response(), "paginate": inputParamStruct.PaginateStruct, "order": inputParamStruct.SortStruct})
}

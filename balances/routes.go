package balances

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testGoProject/users"
)

func UserBalance(router *gin.RouterGroup) {
	router.GET("/:userId", getUserBalance)
}

func UserChangeBalance(router *gin.RouterGroup) {
	router.POST("/:userId/add", addBalance)
	router.POST("/:userId/sub", subBalance)
	router.POST("/:userId/sendTo/:receiverId", exchange)
}

type CurrenciesStruct struct {
	Rates map[string]float64 `json:"rates"`
}

func getUserBalance(context *gin.Context) {
	var inputParam GetBalanceParamStruct

	err := inputParam.BindGetBalanceParams(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	response := new(BalanceResponse)

	err = GetBalanceByUser(inputParam, response)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": users.UserNotFoundMessage})
		return
	}

	context.JSON(http.StatusOK, response)
}

func subBalance(context *gin.Context) {
	var inputParam ChangeParamStruct

	err := inputParam.BindChangeParams(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	message, httpStatus, err := Sub(inputParam)

	if err != nil {
		context.JSON(httpStatus, gin.H{"message": message})
		return
	}

	context.JSON(httpStatus, gin.H{"message": message})
	return
}

func addBalance(context *gin.Context) {
	var inputParam ChangeParamStruct

	err := inputParam.BindChangeParams(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	message, httpStatus, err := Add(inputParam)

	if err != nil {
		context.JSON(httpStatus, gin.H{"message": message})
		return
	}

	context.JSON(httpStatus, gin.H{"message": message})
	return
}

func exchange(context *gin.Context) {
	var inputParam ExchangeParamStruct

	err := inputParam.BindExchangeParams(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	message, httpStatus, err := ExchangeBetweenUsers(inputParam)

	if err != nil {
		context.JSON(httpStatus, gin.H{"message": message})
		return
	}

	context.JSON(httpStatus, gin.H{"message": message})
	return
}

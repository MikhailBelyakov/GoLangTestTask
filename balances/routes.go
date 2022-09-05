package balances

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testProject/users"
)

func BalanceRoutes(router *gin.RouterGroup, controller BalanceController) {
	router.GET("/:userID", controller.GetUserBalance)
	router.POST("/:userID/sub", controller.SubBalance)
	router.POST("/:userID/add", controller.AddBalance)
	router.POST("/:userID/sendTo/:receiverID", controller.Exchange)
}

type BalanceController interface {
	GetUserBalance(ctx *gin.Context)
	SubBalance(ctx *gin.Context)
	AddBalance(ctx *gin.Context)
	Exchange(ctx *gin.Context)
}

func NewBalanceController(service BalanceService) BalanceController {
	return &balanceControllerImpl{
		service: service,
	}
}

type balanceControllerImpl struct {
	service BalanceService
}

func (controller balanceControllerImpl) GetUserBalance(ctx *gin.Context) {
	var inputParam GetBalanceParamStruct

	err := inputParam.BindGetBalanceParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	response := new(BalanceResponse)

	err = controller.service.GetBalanceByUser(ctx, &inputParam, response)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": users.UserNotFoundMessage})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (controller balanceControllerImpl) SubBalance(ctx *gin.Context) {
	var inputParam ChangeParamStruct

	err := inputParam.BindChangeParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	message, httpStatus, err := controller.service.Sub(ctx, &inputParam)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"message": message})
		return
	}

	ctx.JSON(httpStatus, gin.H{"message": message})
	return
}

func (controller balanceControllerImpl) AddBalance(ctx *gin.Context) {
	var inputParam ChangeParamStruct

	err := inputParam.BindChangeParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}

	message, httpStatus, err := controller.service.Add(ctx, &inputParam)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"message": message})
		return
	}

	ctx.JSON(httpStatus, gin.H{"message": message})
	return
}

func (controller balanceControllerImpl) Exchange(ctx *gin.Context) {
	var inputParam ExchangeParamStruct

	err := inputParam.BindExchangeParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": incorrectParamMessage})
		return
	}
	message, httpStatus, err := controller.service.ExchangeBetweenUsers(ctx, &inputParam)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"message": message})
		return
	}

	ctx.JSON(httpStatus, gin.H{"message": message})
	return
}

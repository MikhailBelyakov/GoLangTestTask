package balances

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testProject/common"
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
	var err common.HttpError

	err = inputParam.BindGetBalanceParams(ctx)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.HttpCode()})
		return
	}

	response := new(BalanceResponse)

	err = controller.service.GetBalanceByUser(ctx, &inputParam, response)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (controller balanceControllerImpl) SubBalance(ctx *gin.Context) {
	var inputParam ChangeParamStruct
	var err common.HttpError

	err = inputParam.BindChangeParams(ctx)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	err = controller.service.Sub(ctx, &inputParam)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": subSuccessText})
	return
}

func (controller balanceControllerImpl) AddBalance(ctx *gin.Context) {
	var inputParam ChangeParamStruct
	var err common.HttpError

	err = inputParam.BindChangeParams(ctx)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	err = controller.service.Add(ctx, &inputParam)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": addSuccessText})
	return
}

func (controller balanceControllerImpl) Exchange(ctx *gin.Context) {
	var inputParam ExchangeParamStruct
	var err common.HttpError

	err = inputParam.BindExchangeParams(ctx)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}
	err = controller.service.ExchangeBetweenUsers(ctx, &inputParam)

	if err != nil {
		ctx.JSON(err.HttpCode(), gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": exchangeSuccessText})
	return
}

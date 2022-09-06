package balances

import (
	"github.com/gin-gonic/gin"
)

func BalanceRoutes(router *gin.RouterGroup, controller BalanceController) {
	router.GET("/:userID", controller.GetUserBalance)
	router.POST("/:userID/sub", controller.SubBalance)
	router.POST("/:userID/add", controller.AddBalance)
	router.POST("/:userID/sendTo/:receiverID", controller.Exchange)
}

package transactions

import (
	"github.com/gin-gonic/gin"
)

func UserTransaction(router *gin.RouterGroup, controller TransactionController) {
	router.GET("/:userID", controller.GetUserTransaction)
}

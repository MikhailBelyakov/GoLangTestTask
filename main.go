package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testGoProject/balances"
	"testGoProject/common"
	"testGoProject/transactions"
	"testGoProject/users"
)

func Migrate() {
	users.AutoMigrate()
	balances.AutoMigrate()
	transactions.AutoMigrate()
}

func main() {
	db := common.Init()
	Migrate()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("/api")
	balances.UserBalance(v1.Group("/balances"))
	balances.UserChangeBalance(v1.Group("/balances"))
	transactions.UserTransaction(v1.Group("/transactions"))

	tx1 := db.Begin()
	userA := users.UserModel{
		Username: "testuser",
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)

	r.Run() // listen and serve on 0.0.0.0:8080
}

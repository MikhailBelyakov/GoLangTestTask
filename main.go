package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"sync"
	"testProject/balances"
	"testProject/common"
	"testProject/transactions"
	"testProject/users"
)

func Migrate() {

	var err error

	err = users.AutoMigrate()
	err = balances.AutoMigrate()
	err = transactions.AutoMigrate()

	if err != nil {
		log.Fatal("migration error", err)
	}
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := common.Init()
	Migrate()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("/api")

	balanceRepo, err := balances.NewBalanceRepository(db)

	if err != nil {
		log.Fatal(err.Error())
	}

	transactionRepo, err := transactions.NewTransactionRepository(db)

	if err != nil {
		log.Fatal(err.Error())
	}

	mu := new(sync.Mutex)

	balanceService := balances.NewBalanceService(mu, balanceRepo, transactionRepo)
	transactionService := transactions.NewTransactionService(transactionRepo)

	transactionController := transactions.NewTransactionController(transactionService)
	balanceController := balances.NewBalanceController(balanceService)

	balances.BalanceRoutes(v1.Group("/balances"), balanceController)

	transactions.UserTransaction(v1.Group("/transactions"), transactionController)

	tx1 := db.Begin()
	userA := users.UserModel{
		Username: "testuser",
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)

	r.Run() // listen and serve on 0.0.0.0:8080
}

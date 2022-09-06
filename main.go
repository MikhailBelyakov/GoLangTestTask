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
	if err != nil {
		log.Fatal("migration error", err)
	}

	err = balances.AutoMigrate()
	if err != nil {
		log.Fatal("migration error", err)
	}

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
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	v1 := r.Group("/api")

	balanceRepo := balances.NewBalanceRepository(db)
	transactionRepo := transactions.NewTransactionRepository(db)
	userRepo := users.NewUserRepository(db)

	mu := new(sync.RWMutex)

	balanceService := balances.NewBalanceService(mu, balanceRepo, transactionRepo, userRepo)
	transactionService := transactions.NewTransactionService(transactionRepo)

	balanceController := balances.NewBalanceController(balanceService)
	balances.BalanceRoutes(v1.Group("/balances"), balanceController)

	transactionController := transactions.NewTransactionController(transactionService)
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

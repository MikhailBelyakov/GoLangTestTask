package balances

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"testProject/common"
	"testProject/transactions"
	"testProject/users"
	"time"
)

const exchangeUrl = "https://www.cbr-xml-daily.ru/latest.js"
const defaultVal = "RUB"

type BalanceService interface {
	Sub(ctx *gin.Context, paramStruct *ChangeParamStruct) (string, int, error)
	Add(ctx *gin.Context, paramStruct *ChangeParamStruct) (string, int, error)
	ExchangeBetweenUsers(ctx *gin.Context, inputParam *ExchangeParamStruct) (string, int, error)
	GetBalanceByUser(ctx *gin.Context, inputParam *GetBalanceParamStruct, responseStruct *BalanceResponse) error
	/*getBalance(ctx *gin.Context, userID uint, model *BalanceModel) error
	getCurrency(ctx *gin.Context, currency string, currencyChanel chan<- float64)*/
}

func NewBalanceService(mu *sync.Mutex, balanceRepo BalanceRepository, transactionRepo transactions.TransactionRepository) BalanceService {
	return &balanceServiceImpl{
		mu:              mu,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}
}

type balanceServiceImpl struct {
	mu              *sync.Mutex
	balanceRepo     BalanceRepository
	transactionRepo transactions.TransactionRepository
}

func (service *balanceServiceImpl) Sub(ctx *gin.Context, paramStruct *ChangeParamStruct) (string, int, error) {
	defer service.mu.Unlock()
	service.mu.Lock()

	var balanceModel = new(BalanceModel)

	err := service.getBalance(ctx, paramStruct.userID, balanceModel)

	if err != nil {
		return users.UserNotFoundMessage, http.StatusNotFound, err
	}

	if uint32(paramStruct.amount*100) > balanceModel.Amount {
		return notManyMessage, http.StatusBadRequest, errors.New(notManyMessage)
	}

	newBalance := balanceModel.Amount - uint32(paramStruct.amount*100)
	balanceModel.Amount = newBalance

	tx := common.DB.Begin()

	if err = service.balanceRepo.UpdateBalance(ctx, balanceModel); err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}

	newTransactionModel := transactions.TransactionModel{
		UserID: paramStruct.userID,
		Amount: uint32(paramStruct.amount * 100),
		TypeID: transactions.SubTypeTransaction,
		Date:   time.Now(),
	}

	if err = service.transactionRepo.CreateTransaction(ctx, &newTransactionModel); err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}
	tx.Commit()
	return subSuccessText, http.StatusOK, nil
}

func (service *balanceServiceImpl) Add(ctx *gin.Context, paramStruct *ChangeParamStruct) (string, int, error) {
	defer service.mu.Unlock()
	service.mu.Lock()

	var balanceModel = new(BalanceModel)

	err := service.getBalance(ctx, paramStruct.userID, balanceModel)

	if err != nil {
		return users.UserNotFoundMessage, http.StatusNotFound, err
	}

	newBalance := balanceModel.Amount + uint32(paramStruct.amount*100)
	balanceModel.Amount = newBalance

	tx := common.DB.Begin()

	if err = service.balanceRepo.UpdateBalance(ctx, balanceModel); err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}

	newTransactionModel := transactions.TransactionModel{
		UserID: paramStruct.userID,
		Amount: uint32(paramStruct.amount * 100),
		TypeID: transactions.AddTypeTransaction,
		Date:   time.Now(),
	}

	if err = service.transactionRepo.CreateTransaction(ctx, &newTransactionModel); err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}
	tx.Commit()
	return addSuccessText, http.StatusOK, nil
}

func (service *balanceServiceImpl) ExchangeBetweenUsers(ctx *gin.Context, inputParam *ExchangeParamStruct) (string, int, error) {
	defer service.mu.Unlock()
	service.mu.Lock()

	var senderBalanceModel, receiverBalanceModel = new(BalanceModel), new(BalanceModel)
	var err error

	if inputParam.receiverID == inputParam.senderID {
		return selfSendErrorMessage, http.StatusBadRequest, errors.New(selfSendErrorMessage)
	}

	err = service.getBalance(ctx, inputParam.senderID, senderBalanceModel)

	if err != nil {
		return notFoundSenderMessage, http.StatusNotFound, err
	}

	err = service.getBalance(ctx, inputParam.receiverID, receiverBalanceModel)

	if err != nil {
		return notFoundReceiverMessage, http.StatusNotFound, err
	}

	if uint32(inputParam.amount*100) > senderBalanceModel.Amount {
		return notManyForSendMessage, http.StatusBadRequest, errors.New(notManyForSendMessage)
	}

	newSenderBalance := senderBalanceModel.Amount - uint32(inputParam.amount*100)
	newReceiverBalance := receiverBalanceModel.Amount + uint32(inputParam.amount*100)

	senderBalanceModel.Amount = newSenderBalance
	receiverBalanceModel.Amount = newReceiverBalance

	tx := common.DB.Begin()

	err = service.balanceRepo.UpdateBalance(ctx, senderBalanceModel)
	if err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}
	err = service.balanceRepo.UpdateBalance(ctx, receiverBalanceModel)
	if err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}

	senderTransactionModel := transactions.TransactionModel{
		UserID:     inputParam.senderID,
		SenderID:   inputParam.senderID,
		ReceiverID: inputParam.receiverID,
		Amount:     uint32(inputParam.amount * 100),
		TypeID:     transactions.SendToTypeTransaction,
		Date:       time.Now(),
	}
	receiverTransactionModel := transactions.TransactionModel{
		UserID:     inputParam.receiverID,
		SenderID:   inputParam.senderID,
		ReceiverID: inputParam.receiverID,
		Amount:     uint32(inputParam.amount * 100),
		TypeID:     transactions.ReceiveFromTypeTransaction,
		Date:       time.Now(),
	}

	err = service.transactionRepo.CreateTransaction(ctx, &senderTransactionModel)
	if err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}
	err = service.transactionRepo.CreateTransaction(ctx, &receiverTransactionModel)
	if err != nil {
		tx.Rollback()
		return err.Error(), http.StatusInternalServerError, err
	}

	tx.Commit()

	return exchangeSuccessText, http.StatusOK, nil
}

func (service *balanceServiceImpl) GetBalanceByUser(ctx *gin.Context, inputParam *GetBalanceParamStruct, responseStruct *BalanceResponse) error {
	var currencyChanel = make(chan float64, 1)
	var balanceModel = new(BalanceModel)

	go service.getCurrency(ctx, inputParam.currency, currencyChanel)

	err := service.getBalance(ctx, inputParam.userID, balanceModel)

	if err != nil {
		return err
	}

	currencyFactor := <-currencyChanel
	actualAmount := (float64(balanceModel.Amount) / 100) * currencyFactor

	formatAmount := fmt.Sprintf("%.2f", actualAmount)

	actualAmount, _ = strconv.ParseFloat(formatAmount, 64)
	responseStruct.Amount = actualAmount
	if currencyFactor != 1 {
		responseStruct.Currency = inputParam.currency
	} else {

		responseStruct.Currency = defaultVal
	}

	return nil
}

type CurrenciesStruct struct {
	Rates map[string]float64 `json:"rates"`
}

func (service *balanceServiceImpl) getCurrency(ctx *gin.Context, currency string, currencyChanel chan<- float64) {

	defer func() {
		if r := recover(); r != nil {
			currencyChanel <- 1
			close(currencyChanel)
			return
		}
		close(currencyChanel)
	}()

	if currency != "" {
		var currenciesStruct CurrenciesStruct

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 300*time.Millisecond)
		defer cancel()

		req, _ := http.NewRequest(http.MethodGet, exchangeUrl, nil)

		client := &http.Client{}
		resp, err := client.Do(req.WithContext(timeoutCtx))

		err = json.NewDecoder(resp.Body).Decode(&currenciesStruct)

		if err != nil {
			currencyChanel <- 1
			return
		}

		result := currenciesStruct.Rates[currency]

		if result != 0 {
			currencyChanel <- result
			return
		}

	}
	currencyChanel <- 1
	return
}

func (service *balanceServiceImpl) getBalance(ctx *gin.Context, userID uint, model *BalanceModel) error {

	userModel, err := users.FindOneUser(ctx, &users.UserModel{ID: userID})

	if err != nil {
		return err
	}

	*model = service.balanceRepo.GetBalance(ctx, userModel.ID)
	return nil
}

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

const currencyUrl = "https://www.cbr-xml-daily.ru/latest.js"
const defaultVal = "RUB"
const currencyConnectTimeOut = 500

type BalanceService interface {
	Sub(ctx *gin.Context, paramStruct *ChangeParamStruct) common.HttpError
	Add(ctx *gin.Context, paramStruct *ChangeParamStruct) common.HttpError
	ExchangeBetweenUsers(ctx *gin.Context, inputParam *ExchangeParamStruct) common.HttpError
	GetBalanceByUser(ctx *gin.Context, inputParam *GetBalanceParamStruct, responseStruct *BalanceResponse) common.HttpError
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

func (service *balanceServiceImpl) Sub(ctx *gin.Context, paramStruct *ChangeParamStruct) common.HttpError {
	service.mu.Lock()
	defer service.mu.Unlock()

	var balanceModel = new(BalanceModel)

	err := service.getBalance(ctx, paramStruct.userID, balanceModel)

	if err != nil {
		return common.NewHttpError(http.StatusNotFound, errors.New(users.UserNotFoundMessage))
	}

	if uint32(paramStruct.amount*100) > balanceModel.Amount {
		return common.NewHttpError(http.StatusBadRequest, errors.New(notManyMessage))
	}

	newBalance := balanceModel.Amount - uint32(paramStruct.amount*100)
	balanceModel.Amount = newBalance

	tx := common.DB.Begin()

	if err = service.balanceRepo.UpdateBalance(ctx, balanceModel); err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}

	newTransactionModel := transactions.TransactionModel{
		UserID: paramStruct.userID,
		Amount: uint32(paramStruct.amount * 100),
		TypeID: transactions.SubTypeTransaction,
		Date:   time.Now(),
	}

	if err = service.transactionRepo.CreateTransaction(ctx, &newTransactionModel); err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}
	tx.Commit()
	return nil
}

func (service *balanceServiceImpl) Add(ctx *gin.Context, paramStruct *ChangeParamStruct) common.HttpError {
	service.mu.Lock()
	defer service.mu.Unlock()

	var balanceModel = new(BalanceModel)

	err := service.getBalance(ctx, paramStruct.userID, balanceModel)

	if err != nil {
		return common.NewHttpError(http.StatusNotFound, errors.New(users.UserNotFoundMessage))
	}

	newBalance := balanceModel.Amount + uint32(paramStruct.amount*100)
	balanceModel.Amount = newBalance

	tx := common.DB.Begin()

	if err = service.balanceRepo.UpdateBalance(ctx, balanceModel); err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}

	newTransactionModel := transactions.TransactionModel{
		UserID: paramStruct.userID,
		Amount: uint32(paramStruct.amount * 100),
		TypeID: transactions.AddTypeTransaction,
		Date:   time.Now(),
	}

	if err = service.transactionRepo.CreateTransaction(ctx, &newTransactionModel); err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}
	tx.Commit()
	return nil
}

func (service *balanceServiceImpl) ExchangeBetweenUsers(ctx *gin.Context, inputParam *ExchangeParamStruct) common.HttpError {
	service.mu.Lock()
	defer service.mu.Unlock()

	var senderBalanceModel, receiverBalanceModel = new(BalanceModel), new(BalanceModel)
	var err error

	if inputParam.receiverID == inputParam.senderID {
		return common.NewHttpError(http.StatusBadRequest, errors.New(selfSendErrorMessage))
	}

	err = service.getBalance(ctx, inputParam.senderID, senderBalanceModel)

	if err != nil {
		return common.NewHttpError(http.StatusNotFound, errors.New(notFoundSenderMessage))
	}

	err = service.getBalance(ctx, inputParam.receiverID, receiverBalanceModel)

	if err != nil {

		return common.NewHttpError(http.StatusNotFound, errors.New(notFoundReceiverMessage))
	}

	if uint32(inputParam.amount*100) > senderBalanceModel.Amount {
		return common.NewHttpError(http.StatusBadRequest, errors.New(notManyForSendMessage))
	}

	newSenderBalance := senderBalanceModel.Amount - uint32(inputParam.amount*100)
	newReceiverBalance := receiverBalanceModel.Amount + uint32(inputParam.amount*100)

	senderBalanceModel.Amount = newSenderBalance
	receiverBalanceModel.Amount = newReceiverBalance

	tx := common.DB.Begin()

	err = service.balanceRepo.UpdateBalance(ctx, senderBalanceModel)
	if err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}
	err = service.balanceRepo.UpdateBalance(ctx, receiverBalanceModel)
	if err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
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
		return common.NewHttpError(http.StatusInternalServerError, err)
	}
	err = service.transactionRepo.CreateTransaction(ctx, &receiverTransactionModel)
	if err != nil {
		tx.Rollback()
		return common.NewHttpError(http.StatusInternalServerError, err)
	}

	tx.Commit()

	return nil
}

func (service *balanceServiceImpl) GetBalanceByUser(ctx *gin.Context, inputParam *GetBalanceParamStruct, responseStruct *BalanceResponse) common.HttpError {
	var currencyChanel = make(chan float64, 1)
	var balanceModel = new(BalanceModel)

	go service.getCurrency(ctx, inputParam.currency, currencyChanel)

	err := service.getBalance(ctx, inputParam.userID, balanceModel)

	if err != nil {
		return common.NewHttpError(http.StatusNotFound, errors.New(users.UserNotFoundMessage))
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

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), currencyConnectTimeOut*time.Millisecond)
		defer cancel()

		req, _ := http.NewRequest(http.MethodGet, currencyUrl, nil)

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

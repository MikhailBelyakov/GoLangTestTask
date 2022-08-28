package balances

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testGoProject/transactions"
	"testGoProject/users"
	"time"
)

const exchangeUrl = "https://www.cbr-xml-daily.ru/latest.js"

func Sub(paramStruct ChangeParamStruct) (string, int, error) {
	var balanceModel = new(BalanceModel)

	err := getBalance(paramStruct.userId, balanceModel)

	if err != nil {
		return userNotFoundMessage, http.StatusNotFound, err
	}

	if uint32(paramStruct.amount*100) > balanceModel.Amount {
		return notManyMessage, http.StatusBadRequest, err
	}

	newBalance := balanceModel.Amount - uint32(paramStruct.amount*100)

	if err = balanceModel.changeBalance(newBalance); err != nil {
		return err.Error(), http.StatusInternalServerError, err
	}

	model := transactions.TransactionModel{
		UserID:    paramStruct.userId,
		Amount:    uint32(paramStruct.amount * 100),
		Operation: subOperationText,
		Date:      time.Now(),
	}

	if err = transactions.CreateTransaction(model); err != nil {
		return err.Error(), http.StatusInternalServerError, err
	}

	return subSuccessText, http.StatusOK, nil
}

func Add(paramStruct ChangeParamStruct) (string, int, error) {
	var balanceModel = new(BalanceModel)

	err := getBalance(paramStruct.userId, balanceModel)

	if err != nil {
		return userNotFoundMessage, http.StatusNotFound, err
	}

	newBalance := balanceModel.Amount + uint32(paramStruct.amount*100)

	if err = balanceModel.changeBalance(newBalance); err != nil {
		return err.Error(), http.StatusInternalServerError, err
	}

	model := transactions.TransactionModel{
		UserID:    paramStruct.userId,
		Amount:    uint32(paramStruct.amount * 100),
		Operation: addOperationText,
		Date:      time.Now(),
	}

	if err = transactions.CreateTransaction(model); err != nil {
		return err.Error(), http.StatusInternalServerError, err
	}

	return addSuccessText, http.StatusOK, nil
}

func ExchangeBetweenUsers(inputParam ExchangeParamStruct) (string, int, error) {
	var senderBalanceModel, receiverBalanceModel = new(BalanceModel), new(BalanceModel)
	var err error

	err = getBalance(inputParam.senderId, senderBalanceModel)

	if err != nil {
		return notFoundSenderMessage, http.StatusNotFound, err
	}

	err = getBalance(inputParam.receiverId, receiverBalanceModel)

	if err != nil {
		return notFoundReceiverMessage, http.StatusNotFound, err
	}

	if uint32(inputParam.amount*100) > senderBalanceModel.Amount {
		return notManyForSendMessage, http.StatusBadRequest, err
	}

	newSenderBalance := senderBalanceModel.Amount - uint32(inputParam.amount*100)
	newReceiverBalance := receiverBalanceModel.Amount + uint32(inputParam.amount*100)

	err = senderBalanceModel.changeBalance(newSenderBalance)
	err = receiverBalanceModel.changeBalance(newReceiverBalance)

	senderTransactionModel := transactions.TransactionModel{
		UserID:     inputParam.senderId,
		SenderId:   inputParam.senderId,
		ReceiverId: inputParam.receiverId,
		Amount:     uint32(inputParam.amount * 100),
		Operation:  sendToOperationText,
		Date:       time.Now(),
	}
	receiverTransactionModel := transactions.TransactionModel{
		UserID:     inputParam.receiverId,
		SenderId:   inputParam.senderId,
		ReceiverId: inputParam.receiverId,
		Amount:     uint32(inputParam.amount * 100),
		Operation:  receiveFromOperationText,
		Date:       time.Now(),
	}

	err = transactions.CreateTransaction(senderTransactionModel)
	err = transactions.CreateTransaction(receiverTransactionModel)

	if err != nil {
		return err.Error(), http.StatusInternalServerError, err
	}

	return exchangeSuccessText, http.StatusOK, nil
}

func GetBalanceByUser(inputParam GetBalanceParamStruct, responseStruct *BalanceResponse) error {
	var currencyChanel = make(chan float64, 1)
	var balanceModel = new(BalanceModel)

	go getCurrency(inputParam.currency, currencyChanel)

	err := getBalance(inputParam.userId, balanceModel)

	if err != nil {
		return err
	}

	actualAmount := (float64(balanceModel.Amount) / 100) * <-currencyChanel

	formatAmount := fmt.Sprintf("%.2f", actualAmount)

	actualAmount, _ = strconv.ParseFloat(formatAmount, 64)
	responseStruct.Amount = actualAmount

	return nil
}

func getBalance(userId uint, model *BalanceModel) error {
	userModel, err := users.FindOneUser(&users.UserModel{ID: userId})

	if err != nil {
		return err
	}

	*model = GetBalance(userModel)
	return nil
}

func getCurrency(currency string, currencyChanel chan<- float64) {
	defer close(currencyChanel)
	if currency != "" {
		var currenciesStruct CurrenciesStruct
		resp, _ := http.Get(exchangeUrl)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				currencyChanel <- 1
			}
		}(resp.Body)

		err := json.NewDecoder(resp.Body).Decode(&currenciesStruct)

		if err != nil {
			currencyChanel <- 1
		}

		result := currenciesStruct.Rates[currency]
		if result != 0 {
			currencyChanel <- result
		}

	}
	currencyChanel <- 1
}

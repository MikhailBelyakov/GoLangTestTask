package transactions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type TransactionsStruct struct {
	C                *gin.Context
	TransactionModel []TransactionModel
}

type TransactionStruct struct {
	C *gin.Context
	TransactionModel
}

type TransactionResponse struct {
	Amount     float64   `json:"amount"`
	Operation  string    `json:"operation"`
	Date       time.Time `json:"date"`
	ReceiverId uint      `json:"receiver_id"`
	SenderId   uint      `json:"sender_id"`
}

func (s *TransactionsStruct) Response() []TransactionResponse {
	var response []TransactionResponse
	for _, transaction := range s.TransactionModel {
		serializer := TransactionStruct{s.C, transaction}
		response = append(response, serializer.Response())
	}
	return response
}

func (s *TransactionStruct) Response() TransactionResponse {

	formatAmount := fmt.Sprintf("%.2f", float64(s.Amount)/100)
	actualAmount, _ := strconv.ParseFloat(formatAmount, 64)

	response := TransactionResponse{
		Amount:     actualAmount,
		Operation:  s.Operation,
		ReceiverId: s.ReceiverId,
		SenderId:   s.SenderId,
		Date:       s.Date,
	}
	return response
}

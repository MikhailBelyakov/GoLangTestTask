package transactions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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
	Amount     float64                 `json:"amount"`
	Date       string                  `json:"date"`
	ReceiverID uint                    `json:"receiverID"`
	SenderID   uint                    `json:"senderID"`
	Type       TransactionTypeResponse `json:"type"`
}

type TransactionTypeResponse struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
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

	operationText, _ := GetLabel(s.TypeID)

	response := TransactionResponse{
		Amount:     actualAmount,
		ReceiverID: s.ReceiverID,
		SenderID:   s.SenderID,
		Date:       s.Date.Format("01.02.2006 03:04"),
		Type: TransactionTypeResponse{
			ID:    s.TypeID,
			Value: operationText,
		},
	}
	return response
}

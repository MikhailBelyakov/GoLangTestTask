package transactions

import "errors"

const (
	AddTypeTransaction = iota
	SubTypeTransaction
	SendToTypeTransaction
	ReceiveFromTypeTransaction
)

func Labels() map[int]string {
	return map[int]string{
		AddTypeTransaction:         addOperationText,
		SubTypeTransaction:         subOperationText,
		SendToTypeTransaction:      sendToOperationText,
		ReceiveFromTypeTransaction: receiveFromOperationText,
	}
}

func GetLabel(typeID int) (string, error) {
	if value, inMap := Labels()[typeID]; inMap {
		return value, nil
	}
	return "", errors.New("enum val not found")
}

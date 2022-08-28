package transactions

import (
	"testGoProject/users"
)

func GetTransactionsByUser(paramStruct GetTransactionParamStruct) ([]TransactionModel, error) {
	userModel, err := users.FindOneUser(&users.UserModel{ID: paramStruct.userId})

	if err != nil {
		return []TransactionModel{}, err
	}

	return FindTransactionsByUser(userModel), nil
}

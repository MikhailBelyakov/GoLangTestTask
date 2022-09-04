package balances

type BalanceResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

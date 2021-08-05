package payments

import "math/big"

type AddRequest struct {
	Amount big.Int `json:"amount"`
}
type AddResponse struct {
	Userid         string  `json:"user_id"`
	Amount         big.Int `json:"amount"`
	Type           string  `json:"type"`
	CurrentBalance big.Int `json:"current_balance"`
	Message        string  `json:"message"`
}
type WithdrawRequest struct {
	Amount big.Int `json:"amount"`
}
type WithdrawResponse struct {
	Userid         string  `json:"user_id"`
	Amount         big.Int `json:"amount"`
	Type           string  `json:"type"`
	CurrentBalance big.Int `json:"current_balance"`
	Message        string  `json:"message"`
}

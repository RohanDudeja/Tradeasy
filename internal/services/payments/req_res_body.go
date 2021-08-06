package payments

type AddRequest struct {
	Amount int64 `json:"amount"`
}
type AddResponse struct {
	Userid         string  `json:"user_id"`
	Amount         int64 `json:"amount"`
	Type           string  `json:"type"`
	CurrentBalance int64 `json:"current_balance"`
	Message        string  `json:"message"`
}
type WithdrawRequest struct {
	Amount int64 `json:"amount"`
}
type WithdrawResponse struct {
	Userid         string  `json:"user_id"`
	Amount         int64 `json:"amount"`
	Type           string  `json:"type"`
	CurrentBalance int64 `json:"current_balance"`
	Message        string  `json:"message"`
}

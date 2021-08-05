package payments

type AddRequest struct {
	Userid string `json:"user_id"`
	Amount int    `json:"amount"`
}
type AddResponse struct {
	Userid         string `json:"user_id"`
	Amount         int    `json:"amount"`
	Type           string `json:"type"`
	CurrentBalance int    `json:"current_balance"`
	Message        string `json:"message"`
}
type WithdrawRequest struct {
	Userid string `json:"user_id"`
	Amount int    `json:"amount"`
}
type WithdrawResponse struct {
	Userid         string `json:"user_id"`
	Amount         int    `json:"amount"`
	Type           string `json:"type"`
	CurrentBalance int    `json:"current_balance"`
	Message        string `json:"message"`
}
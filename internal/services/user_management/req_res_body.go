package user_management

type SignUpRequest struct {
	EmailId  string `json:"emailId"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

type UserDetailsRequest struct {
	PanCardNo string `json:"pan_card_no"`
	BankAccNo string `json:"bank_acc_no"`
}

type UserDetailsResponse struct {
	TradingAccId string `json:"trading_acc_id"`
	Balance      int64  `json:"balance"`
	Message      string `json:"message"`
}
type SignInRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Message string `json:"message"`
}
type ForgetPasswordRequest struct {
	UserId  string `json:"user_id"`
	EmailId string `json:"emailId"`
}

type ForgetPasswordResponse struct {
	Otp string `json:"otp"`
}

type VerifyRequest struct {
	UserId      string `json:"user_id"`
	EmailId     string `json:"emailId"`
	NewPassword string `json:"new_password"`
	Otp         string `json:"otp"`
}

type VerifyResponse struct {
	UserId      string `json:"user_id"`
	NewPassword string `json:"new_password"`
	Message     string `json:"message"`
}

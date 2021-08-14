package payments

type AddRequest struct {
	Amount int64 `json:"amount"`
}
type AddResponse struct {
	Message     string `json:"message"`
	PaymentLink string `json:"payment_link"`
}
type WithdrawRequest struct {
	Amount int64 `json:"amount"`
}
type WithdrawResponse struct {
	Amount         int64  `json:"amount"`
	Type           string `json:"type"`
	CurrentBalance int64  `json:"current_balance"`
	Message        string `json:"message"`
}

type RazorpayRequest struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	AcceptPartial  bool   `json:"accept_partial"`
	CallbackURL    string `json:"callback_url"`
	CallbackMethod string `json:"callback_method"`
}

type RazorpayResponse struct {
	AcceptPartial         bool          `json:"accept_partial"`
	Amount                int64         `json:"amount"`
	AmountPaid            int64         `json:"amount_paid"`
	CallbackMethod        string        `json:"callback_method"`
	CallbackURL           string        `json:"callback_url"`
	CancelledAt           int           `json:"cancelled_at"`
	CreatedAt             int           `json:"created_at"`
	Currency              string        `json:"currency"`
	Description           string        `json:"description"`
	ExpireBy              int           `json:"expire_by"`
	ExpiredAt             int           `json:"expired_at"`
	FirstMinPartialAmount int           `json:"first_min_partial_amount"`
	ID                    string        `json:"id"`
	Payments              interface{}   `json:"payments"`
	ReferenceID           string        `json:"reference_id"`
	ReminderEnable        bool          `json:"reminder_enable"`
	Reminders             []interface{} `json:"reminders"`
	ShortURL              string        `json:"short_url"`
	Status                string        `json:"status"`
	UpdatedAt             int           `json:"updated_at"`
	UpiLink               bool          `json:"upi_link"`
	UserID                string        `json:"user_id"`
}

type CallbackParamRequest struct {
	RazorpayPaymentID              string `json:"razorpay_payment_id" form:"razorpay_payment_id"`
	RazorpayPaymentLinkID          string `json:"razorpay_payment_link_id" form:"razorpay_payment_link_id"`
	RazorpayPaymentLinkReferenceID string `json:"razorpay_payment_link_reference_id" form:"razorpay_payment_link_reference_id"`
	RazorpayPaymentLinkStatus      string `json:"razorpay_payment_link_status" form:"razorpay_payment_link_status"`
	RazorpaySignature              string `json:"razorpay_signature" form:"razorpay_signature"`
}

type CallbackResponse struct {
	Balance int64  `json:"balance"`
	Status  string `json:"status"`
}

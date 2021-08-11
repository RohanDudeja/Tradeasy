package payments

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Withdraw  = "withdraw"
	Add       = "add"
	Pending   = "pending"
	Success   = "success"
	RzpKey    = "rzp_test_Oqf3eW39O728uq"
	RzpSecret = "D8vxJixLkoCgWryf1YoCCKp7"
)

func AddAmount(addReq AddRequest, Userid string) (addRes AddResponse, err error) {
	addAmount := addReq.Amount
	var (
		tradingAcc  model.TradingAccount
		addResponse AddResponse
		razorpayRes RazorpayResponse
	)
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return addRes, errors.New("trading account not found")
	}
	razorRequest := RazorpayRequest{
		Amount:         addAmount,
		CallbackURL:    "http://localhost:8080/payments/payment_status",
		CallbackMethod: "get",
		AcceptPartial:  false,
		Currency:       "INR"}
	jsonReq, err := json.Marshal(razorRequest)
	if err != nil {
		return addRes, errors.New("Invalid Request")
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/payment_links", bytes.NewBuffer(jsonReq))
	if err != nil {
		return addRes, errors.New("failed to initiate the payment link ")
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(RzpKey, RzpSecret)
	response, err := client.Do(req)
	if err != nil {
		return addRes, errors.New("invalid key or secret")
	}
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bodyBytes, &razorpayRes)
	if err != nil {
		return addRes, errors.New("failed to read the response")
	}
	pay := model.Payments{
		UserId:         Userid,
		Amount:         addAmount / 100,
		RazorpayLink:   razorpayRes.ShortURL,
		RazorpayLinkId: razorpayRes.ID,
		Status:         Pending,
		PaymentType:    Add,
		CreatedAt:      time.Now()}
	if err = config.DB.Table("payments").Create(&pay).Error; err != nil {
		return addRes, errors.New("payment failed")
	}
	addResponse.Message = "Click the payment link for payment"
	addResponse.PaymentLink = razorpayRes.ShortURL
	return addResponse, err

}
func WithdrawAmount(withdrawReq WithdrawRequest, Userid string) (withdrawRes WithdrawResponse, err error) {
	withdrawAmount := withdrawReq.Amount
	var (
		tradingAcc       model.TradingAccount
		withdrawResponse WithdrawResponse
	)
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return withdrawRes, errors.New("trading account not found")
	}
	if tradingAcc.Balance < withdrawAmount {
		return withdrawRes, errors.New("no sufficient balance in trading account")
	} else {
		tradingAcc.Balance -= withdrawAmount
	}
	if err = config.DB.Table("trading_account").
		Where("user_id = ?", Userid).UpdateColumn("balance", tradingAcc.Balance).Error; err != nil {
		return withdrawRes, errors.New("unable to update the balance")
	}
	pay := model.Payments{
		UserId:      Userid,
		Amount:      withdrawAmount,
		Status:      Success,
		PaymentType: Withdraw,
		CreatedAt:   time.Now()}
	if err = config.DB.Create(&pay).Error; err != nil {
		return withdrawRes, errors.New("payment  failed")
	}
	withdrawResponse.Userid = pay.UserId
	withdrawResponse.Amount = pay.Amount
	withdrawResponse.Type = Withdraw
	withdrawResponse.CurrentBalance = tradingAcc.Balance
	withdrawResponse.Message = "Process Successful"
	return withdrawResponse, err
}

func Callback(request CallbackParamRequest) (callbackRes CallbackResponse, err error) {
	var (
		payments         model.Payments
		callbackResponse CallbackResponse
		tradingAcc       model.TradingAccount
	)
	if err = config.DB.Table("payments").Where("razorpay_link_id=?", request.RazorpayPaymentLinkID).First(&payments).Error; err != nil {
		return callbackRes, errors.New("invalid payment link")
	}
	if err = config.DB.Table("trading_account").Where("user_id=?", payments.UserId).First(&tradingAcc).Error; err != nil {
		return callbackRes, errors.New("trading account not found")
	}
	finalBalance := tradingAcc.Balance + payments.Amount
	if err = config.DB.Table("trading_account").Update("balance", finalBalance).Error; err != nil {
		return callbackRes, errors.New("update balance failed")
	}
	if err = config.DB.Table("payments").
		Where("user_id = ? AND razorpay_link_id=?", payments.UserId, request.RazorpayPaymentLinkID).
		Update("status", Success).Error; err != nil {
		return callbackRes, errors.New("status not updated")
	}
	callbackResponse.Balance = finalBalance
	callbackResponse.Status = "success"
	return callbackResponse, nil
}

package payments

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func AddAmount(addReq AddRequest, Userid string) (addRes *AddResponse, err error) {
	addAmount := addReq.Amount
	var tradingAcc model.TradingAccount
	var addResponse AddResponse
	var razorpayRes RazorpayResponse
	var payments model.Payments
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	razorRequest := RazorpayRequest{Amount: addAmount, CallbackURL: "http://localhost:8080/payments/payment_status", CallbackMethod: "get", AcceptPartial: false}
	jsonReq, err := json.Marshal(razorRequest)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.razorpay.com/v1/payment_links", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &razorpayRes)
	if err != nil {
		return nil, err
	}
	pay := model.Payments{UserId: Userid, Amount: addAmount, RazorpayLink: razorpayRes.ShortURL, RazorpayLinkId: razorpayRes.ID, PaymentType: "add", CreatedAt: time.Now()}
	if err = config.DB.Create(&pay).Error; err != nil {
		return nil, err
	}
	addResponse.Userid = payments.UserId
	addResponse.Amount = payments.Amount
	addResponse.Type = "add"
	addResponse.PaymentLink = razorpayRes.ShortURL
	return &addResponse, err

}
func WithdrawAmount(withdrawReq WithdrawRequest, Userid string) (withdrawRes *WithdrawResponse, err error) {
	withdrawAmount := withdrawReq.Amount
	var tradingAcc model.TradingAccount
	var withdrawResponse WithdrawResponse
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	if tradingAcc.Balance < withdrawAmount {
		return nil, err
	} else {
		tradingAcc.Balance -= withdrawAmount
	}
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).UpdateColumn("balance", tradingAcc.Balance).Error; err != nil {
		return nil, err
	}
	pay := model.Payments{UserId: Userid, Amount: withdrawAmount, CurrentBalance: tradingAcc.Balance, PaymentType: "withdraw", CreatedAt: time.Now()}
	if err = config.DB.Create(&pay).Error; err != nil {
		return nil, err
	}
	withdrawResponse.Userid = pay.UserId
	withdrawResponse.Amount = pay.Amount
	withdrawResponse.Type = "Withdraw"
	withdrawResponse.CurrentBalance = pay.CurrentBalance
	withdrawResponse.Message = "Process Successful"
	return &withdrawResponse, err
}

func Callback(razorpayPaymentID string, razorpayPaymentLinkID string) (callbackRes *CallbackResponse, err error) {
	var payments model.Payments
	var callbackResponse CallbackResponse
	var tradingAcc model.TradingAccount
	if err = config.DB.Table("trading_account").Where("razorpay_link_id=? ,razorpay_link=?", razorpayPaymentID, razorpayPaymentLinkID).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	if err = config.DB.Table("payments").Where("razorpay_link_id=? ,razorpay_link=?", razorpayPaymentID, razorpayPaymentLinkID).First(&payments).Error; err != nil {
		return nil, err
	}
	tradingAcc.Balance += payments.Amount
	userId := tradingAcc.UserId
	if err = config.DB.Table("trading_account").Where("user_id = ?", userId).UpdateColumn("balance", tradingAcc.Balance).Error; err != nil {
		return nil, err
	}
	if err = config.DB.Table("payments").Where("user_id = ?, razorpay_link_id=? ,razorpay_link=?", razorpayPaymentID, razorpayPaymentLinkID).UpdateColumn("current_balance", tradingAcc.Balance).Error; err != nil {
		return nil, err
	}
	callbackResponse.Status = "success"
	return &callbackResponse, nil
}

package payments

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func AddAmount(addReq AddRequest, Userid string) (addRes AddResponse, err error) {
	addAmount := addReq.Amount
	var tradingAcc model.TradingAccount
	var addResponse AddResponse
	var razorpayRes RazorpayResponse
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return addRes, err
	}
	username := "rzp_test_Oqf3eW39O728uq"
	password := "D8vxJixLkoCgWryf1YoCCKp7"
	razorRequest := RazorpayRequest{Amount: addAmount, CallbackURL: "http://localhost:8080/payments/payment_status", CallbackMethod: "get", AcceptPartial: false, Currency: "INR"}
	jsonReq, err := json.Marshal(razorRequest)
	if err != nil {
		return addRes, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/payment_links", bytes.NewBuffer(jsonReq))
	if err != nil {
		return addRes, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(username, password)
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bodyBytes, &razorpayRes)
	if err != nil {
		return addRes, err
	}
	pay := model.Payments{UserId: Userid, Amount: addAmount / 100, RazorpayLink: razorpayRes.ShortURL, RazorpayLinkId: razorpayRes.ID, PaymentType: "add", CreatedAt: time.Now()}
	if err = config.DB.Table("payments").Create(&pay).Error; err != nil {
		return addRes, err
	}
	addResponse.Userid = pay.UserId
	addResponse.Amount = pay.Amount
	addResponse.Type = "add"
	addResponse.PaymentLink = razorpayRes.ShortURL
	return addResponse, err

}
func WithdrawAmount(withdrawReq WithdrawRequest, Userid string) (withdrawRes WithdrawResponse, err error) {
	withdrawAmount := withdrawReq.Amount
	var tradingAcc model.TradingAccount
	var withdrawResponse WithdrawResponse
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return withdrawRes, err
	}
	if tradingAcc.Balance < withdrawAmount {
		return withdrawRes, err
	} else {
		tradingAcc.Balance -= withdrawAmount
	}
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).UpdateColumn("balance", tradingAcc.Balance).Error; err != nil {
		return withdrawRes, err
	}
	pay := model.Payments{UserId: Userid, Amount: withdrawAmount, CurrentBalance: tradingAcc.Balance, PaymentType: "withdraw", CreatedAt: time.Now()}
	if err = config.DB.Create(&pay).Error; err != nil {
		return withdrawRes, err
	}
	withdrawResponse.Userid = pay.UserId
	withdrawResponse.Amount = pay.Amount
	withdrawResponse.Type = "Withdraw"
	withdrawResponse.CurrentBalance = pay.CurrentBalance
	withdrawResponse.Message = "Process Successful"
	return withdrawResponse, err
}

func Callback(razorpayPaymentID string, razorpayPaymentLinkID string) (callbackRes CallbackResponse, err error) {
	var payments model.Payments
	var callbackResponse CallbackResponse
	var tradingAcc model.TradingAccount
	if err = config.DB.Table("payments").Where("razorpay_link_id=?", razorpayPaymentLinkID).First(&payments).Error; err != nil {
		return callbackRes, err
	}
	if err = config.DB.Table("trading_account").Where("user_id=?", payments.UserId).First(&tradingAcc).Error; err != nil {
		return callbackRes, err
	}
	finalBalance := tradingAcc.Balance + payments.Amount
	if err = config.DB.Table("trading_account").Update("balance", finalBalance).Error; err != nil {
		return callbackRes, err
	}
	if err = config.DB.Table("payments").Where("user_id = ? AND razorpay_link_id=?", payments.UserId, razorpayPaymentLinkID).Update("current_balance", finalBalance).Error; err != nil {
		return callbackRes, err
	}
	callbackResponse.Status = "success"
	return callbackResponse, nil
}

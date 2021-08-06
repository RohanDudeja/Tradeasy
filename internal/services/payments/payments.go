package payments

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"sync"
	"time"
)

var Mutex sync.Mutex
func AddAmount(addReq AddRequest, Userid string)(addRes *AddResponse, err error){
	Amount1 := addReq.Amount
	var tradingAcc model.TradingAccount
	var addRes1 AddResponse
	Mutex.Lock()
	if err = config.DB.Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	tradingAcc.Balance += Amount1
	Mutex.Unlock()
	pay := model.Payments{UserId: Userid,RazorpayLinkId: "",RazorpayLink: "",Amount: Amount1,CurrentBalance: tradingAcc.Balance,UpdatedAt: time.Now() }
	if err = config.DB.Create(pay).Error; err != nil {
		return nil, err
	}
	addRes1.Userid = pay.UserId
	addRes1.Amount = pay.Amount
	addRes1.Type = "add"
	addRes1.CurrentBalance = pay.CurrentBalance
	addRes1.Message = "Process Successful"
	return &addRes1,err
}
func WithdrawAmount(withdrawReq WithdrawRequest, Userid string )(withdrawRes *WithdrawResponse, err error){
	Amount1 := withdrawReq.Amount
	var tradingAcc model.TradingAccount
	var withdrawRes1 WithdrawResponse
	Mutex.Lock()
	if err = config.DB.Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	if tradingAcc.Balance < Amount1 {
		return nil,err
	} else {
		tradingAcc.Balance -= Amount1}
	Mutex.Unlock()
	pay := model.Payments{UserId: Userid,RazorpayLinkId: "",RazorpayLink: "",Amount: Amount1,CurrentBalance: tradingAcc.Balance,UpdatedAt: time.Now() }
	if err = config.DB.Create(pay).Error; err != nil {
		return nil, err
	}
	withdrawRes1.Userid = pay.UserId
	withdrawRes1.Amount = pay.Amount
	withdrawRes1.Type = "Withdraw"
	withdrawRes1.CurrentBalance = pay.CurrentBalance
	withdrawRes1.Message = "Process Successful"
	return &withdrawRes1,err
}
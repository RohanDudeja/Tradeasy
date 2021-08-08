package payments

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"sync"
	"time"
)

var Mutex sync.Mutex

func AddAmount(addReq AddRequest, Userid string) (addRes *AddResponse, err error) {
	addAmount := addReq.Amount
	var tradingAcc model.TradingAccount
	var addResponse AddResponse
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).First(&tradingAcc).Error; err != nil {
		return nil, err
	}
	tradingAcc.Balance += addAmount
	if err = config.DB.Table("trading_account").Where("user_id = ?", Userid).UpdateColumn("balance", tradingAcc.Balance).Error; err != nil {
		return nil, err
	}
	pay := model.Payments{UserId: Userid, RazorpayLinkId: "", RazorpayLink: "", Amount: addAmount, CurrentBalance: tradingAcc.Balance, UpdatedAt: time.Now()}
	if err = config.DB.Create(pay).Error; err != nil {
		return nil, err
	}
	addResponse.Userid = pay.UserId
	addResponse.Amount = pay.Amount
	addResponse.Type = "add"
	addResponse.CurrentBalance = pay.CurrentBalance
	addResponse.Message = "Process Successful"
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
	pay := model.Payments{UserId: Userid, RazorpayLinkId: "", RazorpayLink: "", Amount: withdrawAmount, CurrentBalance: tradingAcc.Balance, UpdatedAt: time.Now()}
	if err = config.DB.Create(pay).Error; err != nil {
		return nil, err
	}
	withdrawResponse.Userid = pay.UserId
	withdrawResponse.Amount = pay.Amount
	withdrawResponse.Type = "Withdraw"
	withdrawResponse.CurrentBalance = pay.CurrentBalance
	withdrawResponse.Message = "Process Successful"
	return &withdrawResponse, err
}

package reports

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"fmt"
	"time"
)

func DailyPendingOrders(Userid string) (penOrderRes []DailyPendingOrderResponse, err error) {
	var (
		pendingOrders    []model.PendingOrders
		penOrderResponse []DailyPendingOrderResponse
	)
	if err = config.DB.Table("pending_orders").Where("user_id = ? AND deleted_at = NULL ", Userid).Find(&pendingOrders).Error; err != nil {
		return nil, err
	}
	for _, pend := range pendingOrders {
		var pendingOrderResponse DailyPendingOrderResponse
		pendingOrderResponse.Userid = pend.UserId
		pendingOrderResponse.OrderId = pend.OrderId
		pendingOrderResponse.StockName = pend.StockName
		pendingOrderResponse.OrderType = pend.OrderType
		pendingOrderResponse.BookType = pend.BookType
		pendingOrderResponse.LimitPrice = pend.LimitPrice
		pendingOrderResponse.Quantity = pend.Quantity
		pendingOrderResponse.OrderPrice = pend.OrderPrice
		pendingOrderResponse.Status = pend.Status
		penOrderResponse = append(penOrderResponse, pendingOrderResponse)
	}
	return penOrderResponse, nil
}
func Portfolio(Userid string, from string, to string) (portfolioRes []PortfolioResponse, err error) {
	var (
		portfolio         []model.Holdings
		portfolioResponse []PortfolioResponse
	)
	fromTime, err := time.Parse("020106150405", from)
	if err != nil {
		fmt.Println(err.Error())
	}
	toTime, err := time.Parse("020106150405", to)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := config.DB.Table("holdings").Where("Userid = ? AND updated_at BETWEEN ? AND ? ", Userid, fromTime, toTime).Find(&portfolio).Error; err != nil {
		return nil, err
	}
	for _, portf := range portfolio {
		var portResponse PortfolioResponse
		portResponse.Userid = portf.UserId
		portResponse.OrderId = portf.OrderId
		portResponse.StockName = portf.StockName
		portResponse.Quantity = portf.Quantity
		portResponse.BuyPrice = portf.BuyPrice
		portfolioResponse = append(portfolioResponse, portResponse)
	}
	return portfolioResponse, nil
}
func OrdersHistory(Userid string, from string, to string) (ordHisRes []OrderHistoryResponse, err error) {
	var (
		orderHistory   []model.OrderHistory
		holdings       []model.Holdings
		ordHisResponse []OrderHistoryResponse
	)
	fromTime, err := time.Parse("020106150405", from)
	if err != nil {
		fmt.Println(err.Error())
	}
	toTime, err := time.Parse("020106150405", to)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err = config.DB.Table("order_history").Where("Userid = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&orderHistory).Error; err != nil {
		return nil, err
	}
	if err = config.DB.Table("holdings").Where("Userid = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&holdings).Error; err != nil {
		return nil, err
	}
	for _, ordHis := range orderHistory {
		var ordHistoryRes OrderHistoryResponse
		ordHistoryRes.Userid = ordHis.UserId
		ordHistoryRes.OrderId = ordHis.OrderId
		ordHistoryRes.StockName = ordHis.StockName
		ordHistoryRes.Quantity = ordHis.Quantity
		ordHistoryRes.BuyPrice = ordHis.BuyPrice
		ordHistoryRes.SellPrice = ordHis.SellPrice
		ordHisResponse = append(ordHisResponse, ordHistoryRes)
	}
	for _, hold := range holdings {
		var orderHisRes OrderHistoryResponse
		orderHisRes.Userid = hold.UserId
		orderHisRes.OrderId = hold.OrderId
		orderHisRes.StockName = hold.StockName
		orderHisRes.Quantity = hold.Quantity
		orderHisRes.BuyPrice = hold.BuyPrice
		orderHisRes.SellPrice = 0
		ordHisResponse = append(ordHisResponse, orderHisRes)
	}
	return ordHisResponse, nil
}
func ProfitLossHistory(Userid string, from string, to string) (proLosRes []ProfitLossHistoryResponse, err error) {
	var (
		profitLossHistory  []model.OrderHistory
		profitLossResponse []ProfitLossHistoryResponse
	)
	fromTime, err := time.Parse("020106150405", from)
	if err != nil {
		fmt.Println(err.Error())
	}
	toTime, err := time.Parse("020106150405", to)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err = config.DB.Table("order_history").Where("Userid = ?  AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&profitLossHistory).Error; err != nil {
		return nil, err
	}
	for _, profitloss := range profitLossHistory {
		var proLosResponse ProfitLossHistoryResponse
		proLosResponse.Userid = profitloss.UserId
		proLosResponse.OrderId = profitloss.OrderId
		proLosResponse.StockName = profitloss.StockName
		proLosResponse.Quantity = profitloss.Quantity
		proLosResponse.BuyPrice = profitloss.BuyPrice
		proLosResponse.SellPrice = profitloss.SellPrice
		proLosResponse.ProfitLoss = profitloss.Quantity * (profitloss.SellPrice - profitloss.BuyPrice)
		profitLossResponse = append(profitLossResponse, proLosResponse)
	}
	return profitLossResponse, nil
}
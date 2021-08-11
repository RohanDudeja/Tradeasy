package reports

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"errors"
	"time"
)

func DailyPendingOrders(Userid string) (penOrderRes []DailyPendingOrderResponse, err error) {
	var (
		pendingOrders             []model.PendingOrders
		dailyPendingOrderResponse []DailyPendingOrderResponse
	)
	if err = config.DB.Table("pending_orders").Where("user_id = ?", Userid).Find(&pendingOrders).Error; err != nil {
		return nil, errors.New("no pending orders found")
	}
	for _, pendingOrder := range pendingOrders {
		var pendingOrderResponse DailyPendingOrderResponse
		pendingOrderResponse.Userid = pendingOrder.UserId
		pendingOrderResponse.OrderId = pendingOrder.OrderId
		pendingOrderResponse.StockName = pendingOrder.StockName
		pendingOrderResponse.OrderType = pendingOrder.OrderType
		pendingOrderResponse.BookType = pendingOrder.BookType
		pendingOrderResponse.LimitPrice = pendingOrder.LimitPrice
		pendingOrderResponse.Quantity = pendingOrder.Quantity
		pendingOrderResponse.OrderPrice = pendingOrder.OrderPrice
		pendingOrderResponse.Status = pendingOrder.Status
		dailyPendingOrderResponse = append(dailyPendingOrderResponse, pendingOrderResponse)
	}
	return dailyPendingOrderResponse, nil
}
func Portfolio(Userid string, request ReportsParamRequest) (portfolioRes []PortfolioResponse, err error) {
	var (
		holdings          []model.Holdings
		portfolioResponse []PortfolioResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)

	if err := config.DB.Table("holdings").Where("user_id = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&holdings).Error; err != nil {
		return nil, errors.New("no holdings found")
	}
	for _, holding := range holdings {
		var portResponse PortfolioResponse
		portResponse.Userid = holding.UserId
		portResponse.OrderId = holding.OrderId
		portResponse.StockName = holding.StockName
		portResponse.Quantity = holding.Quantity
		portResponse.BuyPrice = holding.BuyPrice
		portfolioResponse = append(portfolioResponse, portResponse)
	}
	return portfolioResponse, nil
}
func OrdersHistory(Userid string, request ReportsParamRequest) (ordHisRes []OrderHistoryResponse, err error) {
	var (
		orderHistory   []model.OrderHistory
		holdings       []model.Holdings
		ordHisResponse []OrderHistoryResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)
	if err = config.DB.Table("order_history").
		Where("user_id = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&orderHistory).Error; err != nil {
		return nil, errors.New("no orders found")
	}
	if err = config.DB.Table("holdings").
		Where("user_id = ? AND created_at BETWEEN ? AND ? ", Userid, fromTime, toTime).
		Group("order_id").Unscoped().Find(&holdings).Error; err != nil {
		return nil, errors.New("no orders found")

	}
	for _, ordHis := range orderHistory {

		var ordHistoryRes OrderHistoryResponse
		ordHistoryRes.Userid = ordHis.UserId
		ordHistoryRes.OrderId = ordHis.OrderId
		ordHistoryRes.StockName = ordHis.StockName
		ordHistoryRes.Quantity = ordHis.Quantity
		ordHistoryRes.BuySellType = "SELL"
		ordHisResponse = append(ordHisResponse, ordHistoryRes)
	}
	for _, hold := range holdings {
		var orderHisRes OrderHistoryResponse
		orderHisRes.Userid = hold.UserId
		orderHisRes.OrderId = hold.OrderId
		var holdingsQuantity int
		if err = config.DB.Raw("SELECT SUM(quantity) FROM holdings WHERE user_id = ? AND order_id ", hold.UserId, hold.OrderId).
			Scan(&holdingsQuantity).Error; err != nil {
			return nil, errors.New("problem fetching quantity")
		}
		var orderHistoryQuantity int
		if err = config.DB.Raw("SELECT SUM(quantity) FROM order_history WHERE user_id = ? AND order_id ", hold.UserId, hold.OrderId).
			Scan(&orderHistoryQuantity).Error; err != nil {
			return nil, errors.New("problem fetching quantity")
		}
		orderHisRes.StockName = hold.StockName
		orderHisRes.Quantity = holdingsQuantity + orderHistoryQuantity
		orderHisRes.BuySellType = "BUY"
		ordHisResponse = append(ordHisResponse, orderHisRes)
	}
	return ordHisResponse, nil

}
func ProfitLossHistory(Userid string, request ReportsParamRequest) (proLosRes []ProfitLossHistoryResponse, err error) {
	var (
		profitLossHistory  []model.OrderHistory
		profitLossResponse []ProfitLossHistoryResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)
	if err = config.DB.Table("order_history").
		Where("user_id = ?  AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).
		Find(&profitLossHistory).Error; err != nil {
		return nil, errors.New("no orders found")
	}
	for _, profitLoss := range profitLossHistory {
		var proLosResponse ProfitLossHistoryResponse
		proLosResponse.Userid = profitLoss.UserId
		proLosResponse.OrderId = profitLoss.OrderId
		proLosResponse.StockName = profitLoss.StockName
		proLosResponse.Quantity = profitLoss.Quantity
		proLosResponse.BuyPrice = profitLoss.BuyPrice
		proLosResponse.SellPrice = profitLoss.SellPrice
		proLosResponse.ProfitLoss = profitLoss.Quantity * (profitLoss.SellPrice - profitLoss.BuyPrice)
		if len(profitLossResponse) == 0 {
			proLosResponse.CumulatedProfit = profitLoss.Quantity * (profitLoss.SellPrice - profitLoss.BuyPrice)
		} else {
			proLosResponse.CumulatedProfit = profitLoss.Quantity*(profitLoss.SellPrice-profitLoss.BuyPrice) + profitLossResponse[len(profitLossResponse)-1].CumulatedProfit
		}
		profitLossResponse = append(profitLossResponse, proLosResponse)
	}
	return profitLossResponse, nil

}

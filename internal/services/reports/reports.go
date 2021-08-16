package reports

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"errors"
	"time"
)

const (
	Buy  = "buy"
	Sell = "sell"
)

func DailyPendingOrders(Userid string) (response []DailyPendingOrderResponse, err error) {
	var (
		pendingOrders             []model.PendingOrders
		dailyPendingOrderResponse []DailyPendingOrderResponse
	)
	if err = database.GetDB().Table("pending_orders").Where("user_id = ?", Userid).Find(&pendingOrders).Error; err != nil {
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
func Portfolio(Userid string, request ReportsParamRequest) (response []PortfolioResponse, err error) {
	var (
		holdings           []model.Holdings
		portfolioResponses []PortfolioResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)

	if err := database.GetDB().Table("holdings").
		Where("user_id = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&holdings).Error; err != nil {
		return nil, errors.New("no holdings found")
	}
	for _, holding := range holdings {
		var portfolioResponse PortfolioResponse
		portfolioResponse.Userid = holding.UserId
		portfolioResponse.OrderId = holding.OrderId
		portfolioResponse.StockName = holding.StockName
		portfolioResponse.Quantity = holding.Quantity
		portfolioResponse.BuyPrice = holding.BuyPrice
		portfolioResponses = append(portfolioResponses, portfolioResponse)
	}
	return portfolioResponses, nil
}
func OrdersHistory(Userid string, request ReportsParamRequest) (response []OrderHistoryResponse, err error) {
	var (
		orderHistories        []model.OrderHistory
		holdings              []model.Holdings
		orderHistoryResponses []OrderHistoryResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)
	if err = database.GetDB().Table("order_history").
		Where("user_id = ? AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).Find(&orderHistories).Error; err != nil {
		return nil, errors.New("no order history found")
	}
	if err = database.GetDB().Table("holdings").
		Where("user_id = ? AND created_at BETWEEN ? AND ? ", Userid, fromTime, toTime).
		Group("order_id").Unscoped().Find(&holdings).Error; err != nil {
		return nil, errors.New("no order history found")

	}
	for _, orderHistory := range orderHistories {

		var orderHistoryResponse OrderHistoryResponse
		orderHistoryResponse.Userid = orderHistory.UserId
		orderHistoryResponse.OrderId = orderHistory.OrderId
		orderHistoryResponse.StockName = orderHistory.StockName
		orderHistoryResponse.Quantity = orderHistory.Quantity
		orderHistoryResponse.BuySellType = Sell
		orderHistoryResponses = append(orderHistoryResponses, orderHistoryResponse)
	}
	for _, holding := range holdings {
		var orderHistoryResponse OrderHistoryResponse
		orderHistoryResponse.Userid = holding.UserId
		orderHistoryResponse.OrderId = holding.OrderId
		type HoldingsQuantity struct {
			TotalQuantity int
		}
		var holdingsQuantity HoldingsQuantity
		if err = database.GetDB().Table("holdings").Select("sum(quantity) as total_quantity").
			Where("user_id=? AND order_id=?", holding.UserId, holding.OrderId).
			Scan(&holdingsQuantity).Error; err != nil {
			return nil, errors.New("error in fetching quantity")
		}

		var orderHistoryQuantity HoldingsQuantity
		if err = database.GetDB().Table("order_history").Select("sum(quantity) as total_quantity").
			Where("user_id=? AND order_id=?", holding.UserId, holding.OrderId).
			Scan(&orderHistoryQuantity).Error; err != nil {
			return nil, errors.New("error in fetching quantity")
		}
		orderHistoryResponse.StockName = holding.StockName
		orderHistoryResponse.Quantity = holdingsQuantity.TotalQuantity + orderHistoryQuantity.TotalQuantity
		orderHistoryResponse.BuySellType = Buy
		orderHistoryResponses = append(orderHistoryResponses, orderHistoryResponse)
	}
	return orderHistoryResponses, nil

}
func ProfitLossHistory(Userid string, request ReportsParamRequest) (response []ProfitLossHistoryResponse, err error) {
	var (
		profitLossHistories []model.OrderHistory
		profitLossResponses []ProfitLossHistoryResponse
	)
	fromTime := time.Unix(int64(request.From), 0)
	toTime := time.Unix(int64(request.To), 0)
	if err = database.GetDB().Table("order_history").
		Where("user_id = ?  AND updated_at BETWEEN ? AND ?", Userid, fromTime, toTime).
		Find(&profitLossHistories).Error; err != nil {
		return nil, errors.New("no order history found for profit loss report")
	}
	for _, profitLossHistory := range profitLossHistories {
		var profitLossResponse ProfitLossHistoryResponse
		profitLossResponse.Userid = profitLossHistory.UserId
		profitLossResponse.OrderId = profitLossHistory.OrderId
		profitLossResponse.StockName = profitLossHistory.StockName
		profitLossResponse.Quantity = profitLossHistory.Quantity
		profitLossResponse.BuyPrice = profitLossHistory.BuyPrice
		profitLossResponse.SellPrice = profitLossHistory.SellPrice
		profitLossResponse.ProfitLoss = profitLossHistory.Quantity * (profitLossHistory.SellPrice - profitLossHistory.BuyPrice)
		if len(profitLossResponses) == 0 {
			profitLossResponse.CumulatedProfit = profitLossHistory.Quantity * (profitLossHistory.SellPrice - profitLossHistory.BuyPrice)
		} else {
			profitLossResponse.CumulatedProfit = profitLossHistory.Quantity*(profitLossHistory.SellPrice-profitLossHistory.BuyPrice) +
				profitLossResponses[len(profitLossResponses)-1].CumulatedProfit
		}
		profitLossResponses = append(profitLossResponses, profitLossResponse)
	}
	return profitLossResponses, nil

}

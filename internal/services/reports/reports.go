package reports

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"fmt"
	"time"
)

func DailyPendingOrders(Userid string) (penOrderRes []DailyPendingOrderResponse,err error) {
	var pendingOrders []model.PendingOrders
	var penOrderRes2 []DailyPendingOrderResponse
	if err = config.DB.Where("user_id = ? AND deleted_at = NULL ", Userid).Find(&pendingOrders).Error; err != nil {
		return nil, err
	}
	for _,pend := range pendingOrders{
		var pendingOrderRes1 DailyPendingOrderResponse
		pendingOrderRes1.Userid = pend.UserId
		pendingOrderRes1.OrderId = pend.OrderId
		pendingOrderRes1.StockName = pend.StockName
		pendingOrderRes1.OrderType = pend.OrderType
		pendingOrderRes1.BookType = pend.BookType
		pendingOrderRes1.LimitPrice= pend.LimitPrice
		pendingOrderRes1.Quantity = pend.Quantity
		pendingOrderRes1.OrderPrice = pend.OrderPrice
		pendingOrderRes1.Status = pend.Status
		penOrderRes2 = append(penOrderRes2,pendingOrderRes1)
	}
	return penOrderRes2,nil
}
func Portfolio(Userid string, from string,to string ) ( portfolioRes []PortfolioResponse, err error){
	var portfolio []model.Holdings
	var portfolioRes2 []PortfolioResponse
	from1,err := time.Parse("020106150405",from)
	if err!= nil {
		fmt.Println(err.Error())
	}
	to1,err := time.Parse("020106150405", to)
	if err!= nil {
		fmt.Println(err.Error())
	}
	if err:= config.DB.Where("Userid = ? AND updated_at BETWEEN ? AND ? ", Userid ,from1 ,to1).Find(&portfolio).Error; err != nil {
		return nil,err
	}
	for _,portf := range portfolio{
		var portfolioRes1 PortfolioResponse
		portfolioRes1.Userid = portf.UserId
		portfolioRes1.OrderId = portf.OrderId
		portfolioRes1.StockName = portf.StockName
		portfolioRes1.Quantity = portf.Quantity
		portfolioRes1.BuyPrice = portf.BuyPrice
		portfolioRes2 = append(portfolioRes2,portfolioRes1)
	}
	return portfolioRes2,nil
}
func OrdersHistory(Userid string,from string,to string) (ordHisRes []OrderHistoryResponse,err error){
	var(report1 []model.OrderHistory
		report2 []model.Holdings)
	var ordHisRes2 []OrderHistoryResponse
	from1,err := time.Parse("020106150405",from)
	if err!= nil {
		fmt.Println(err.Error())
	}
	to1,err := time.Parse("020106150405", to)
	if err!= nil {
		fmt.Println(err.Error())
	}
	if err = config.DB.Where("Userid = ? AND updated_at BETWEEN ? AND ?", Userid,from1,to1).Find(&report1).Error; err != nil {
		return nil,err
	}
	if err = config.DB.Where("Userid = ? AND updated_at BETWEEN ? AND ?", Userid,from1,to1).Find(&report2).Error; err != nil {
		return nil,err
	}
	for _,ordHis := range report1{
		var ordHisRes1 OrderHistoryResponse
		ordHisRes1.Userid = ordHis.UserId
		ordHisRes1.OrderId = ordHis.OrderId
		ordHisRes1.StockName = ordHis.StockName
		ordHisRes1.Quantity = ordHis.Quantity
		ordHisRes1.BuyPrice = ordHis.BuyPrice
		ordHisRes1.SellPrice = ordHis.SellPrice
		ordHisRes2 = append(ordHisRes2,ordHisRes1)
	}
	for _,hold := range report2{
		var ordHisRes3 OrderHistoryResponse
		ordHisRes3.Userid = hold.UserId
		ordHisRes3.OrderId = hold.OrderId
		ordHisRes3.StockName = hold.StockName
		ordHisRes3.Quantity = hold.Quantity
		ordHisRes3.BuyPrice = hold.BuyPrice
		ordHisRes3.SellPrice = 0
		ordHisRes2 = append(ordHisRes2,ordHisRes3)
	}
	return ordHisRes2,nil
}
func ProfitLossHistory( Userid string, from string,to string) (proLosRes []ProfitLossHistoryResponse,err error){
	var profitLossHistory [] model.OrderHistory
	var proLosRes2 [] ProfitLossHistoryResponse
	from1,err := time.Parse("020106150405",from)
	if err!= nil {
		fmt.Println(err.Error())
	}
	to1,err := time.Parse("020106150405", to)
	if err!= nil {
		fmt.Println(err.Error())
	}
	if err = config.DB.Where("Userid = ?  AND updated_at BETWEEN ? AND ?", Userid,from1,to1).Find(&profitLossHistory).Error; err != nil {
		return nil,err
	}
	for _,prolos := range profitLossHistory{
		var proLosRes1 ProfitLossHistoryResponse
		proLosRes1.Userid = prolos.UserId
		proLosRes1.OrderId = prolos.OrderId
		proLosRes1.StockName = prolos.StockName
		proLosRes1.Quantity = prolos.Quantity
		proLosRes1.BuyPrice = prolos.BuyPrice
		proLosRes1.SellPrice= prolos.SellPrice
		proLosRes1.ProfitLoss = prolos.Quantity*(prolos.SellPrice-prolos.BuyPrice)
		proLosRes2 = append(proLosRes2,proLosRes1)
	}
	return proLosRes2,nil
}

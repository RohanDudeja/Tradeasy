package reports

import (
	"Tradeasy/internal/model"
)

func DailyPendingOrders(reports []model.PendingOrders, Userid string) (penOrderRes DailyPendingOrderResponse, err error) {
	return penOrderRes, nil

}
func Portfolio(reports []model.Holdings, Userid string, from string, to string) (portfolioRes PortfolioResponse, err error) {

	return portfolioRes, nil
}
func OrdersHistory(report1 []model.OrderHistory, report2 []model.Holdings, Userid string, from string, to string) (ordHisRes OrderHistoryResponse, err error) {

	return ordHisRes, nil
}
func ProfitLossHistory(reports []model.OrderHistory, Userid string, from string, to string) (proLosRes ProfitLossHistoryResponse, err error) {

	return proLosRes, nil
}

package router

import (
	"Tradeasy/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	trade := r.Group("/pending_orders")
	{

		trade.POST(":Userid/buy", controller.BuyOrder)
		trade.POST(":Userid/sell", controller.SellOrder)
		trade.PATCH(":OrderId/cancel", controller.CancelOrder)
	}

	exchangeBuy := r.Group("/buy_order_book")
	{
		exchangeBuy.POST("buy_order", controller.ExecuteBuyOrder)
		exchangeBuy.DELETE("buy_order/:order_id", controller.DeleteBuyOrder)
	}
	exchangeSell := r.Group("/sell_order_book")
	{
		exchangeSell.POST("sell_order", controller.ExecuteSellOrder)
		exchangeSell.DELETE("sell_order/:order_id", controller.DeleteSellOrder)
	}
	exchangeFetch := r.Group("/order_book")
	{
		exchangeFetch.GET(":stock_name/depth", controller.ViewMarketDepth)
	}

	//websocket := r.Group("/socket")
	//{
	//	websocket.GET("/", webSocket.Home)
	//	websocket.GET("/stocks", webSocket.StockHandler)
	//	websocket.GET("/orders", webSocket.OrderHandler)
	//}
	payments := r.Group("/payments")
	{
		payments.POST(":Userid/addAmount", controller.AddAmount)
		payments.POST(":Userid/withdrawAmount", controller.WithdrawAmount)

	}

	reports := r.Group("/reports")
	{
		reports.GET("pending_orders/:Userid", controller.DailyPendingOrders)
		reports.GET("holdings/:Userid", controller.Portfolio)
		reports.GET("order_history/:Userid", controller.OrdersHistory)
		reports.GET("profit_loss_history/:Userid", controller.ProfitLossHistory)
	}
	return r
}

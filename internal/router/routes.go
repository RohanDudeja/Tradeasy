package router

import (
	"Tradeasy/internal/controller"
	"Tradeasy/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	trade := r.Group("/pending_orders")
	trade.Use(middleware.UserBasicAuth())
	{

		trade.POST(":Userid/buy", controller.BuyOrder)
		trade.POST(":Userid/sell", controller.SellOrder)
		trade.PATCH(":OrderId/cancel", controller.CancelOrder)
	}

	exchangeBuy := r.Group("/buy_order_book")
	exchangeBuy.Use(middleware.ExchangeBasicAuth())
	{
		exchangeBuy.POST("buy_order", controller.ExecuteBuyOrder)
		exchangeBuy.DELETE("buy_order/:order_id", controller.DeleteBuyOrder)
	}
	exchangeSell := r.Group("/sell_order_book")
	exchangeSell.Use(middleware.ExchangeBasicAuth())
	{
		exchangeSell.POST("sell_order", controller.ExecuteSellOrder)
		exchangeSell.DELETE("sell_order/:order_id", controller.DeleteSellOrder)
	}
	exchangeFetch := r.Group("/order_book")
	exchangeFetch.Use(middleware.ExchangeBasicAuth())
	{
		exchangeFetch.GET(":stock_name/depth", controller.ViewMarketDepth)
	}

	//websocket:= r.Group("/socket")
	//{
	//	websocket.GET("/", webSocket.Home)
	//	websocket.GET("/stocks", webSocket.StockHandler)
	//	websocket.GET("/orders", webSocket.OrderHandler)
	//}
	watchlist := r.Group("/user_watchlist")
	watchlist.Use(middleware.ExchangeBasicAuth())
	{
		watchlist.POST("", controller.CreateWatchlist)
		watchlist.POST("/:watchlist_id/add", controller.AddStockEntry)
		watchlist.DELETE("/:watchlist_id", controller.DeleteStockEntry)
		watchlist.PATCH("/sort", controller.SortWatchlist)
	}
	return r
}

package router

import (
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine  {
	r:=gin.Default()
	/*
	trade:=r.Group("/pending_orders")
	{

		trade.POST(":Userid/buy",controller.BuyOrder)
		trade.POST(":Userid/sell",controller.SellOrder)
		trade.PATCH(":OrderId/cancel",controller.CancelOrder)
	}

	exchange := r.Group("/order_book")
	{
		exchange.POST("buy_order", controller.ExecuteBuyOrder)
		exchange.POST("sell_order", controller.ExecuteSellOrder)
		exchange.DELETE("buy_order/:StockName", controller.DeleteBuyOrder)
		exchange.DELETE("sell_order/:StockName", controller.DeleteSellOrder)
		exchange.GET(":StockName/depth", controller.ViewMarketDepth)
	}
	 */
	return r
}

package routes

import (
	"Tradeasy/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetUpOrderRouter() *gin.Engine  {
	r:=gin.Default()
	trade:=r.Group("/pending_orders")
	{
		trade.POST(":Userid/buy",controller.BuyOrder)
		trade.POST(":Userid/sell",controller.SellOrder)
		trade.PATCH(":OrderId/cancel",controller.CancelOrder)
	}
	return r
}

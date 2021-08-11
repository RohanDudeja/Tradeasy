package router

import (
	"Tradeasy/internal/controller"
	webSocket "Tradeasy/internal/controller"
	"Tradeasy/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	trade := r.Group("/pending_orders")
	trade.Use(middleware.UserBasicAuth())
	{

		trade.POST(":user_id/buy", controller.BuyOrder)
		trade.POST(":user_id/sell", controller.SellOrder)
		trade.PATCH(":order_id/cancel", controller.CancelOrder)
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

	websocket := r.Group("/socket")
	{
		websocket.GET("/stocks", webSocket.StockHandler)
		websocket.GET("/orders", webSocket.OrderHandler)
	}

	watchlist := r.Group("/user_watchlist")
	watchlist.Use(middleware.UserBasicAuth())
	{
		watchlist.POST("", controller.CreateWatchlist)
		watchlist.POST("/:watchlist_id/add", controller.AddStockEntry)
		watchlist.DELETE("/:watchlist_id", controller.DeleteStockEntry)
		watchlist.PATCH("/sort", controller.SortWatchlist)
	}

	//userSign no auth needed
	userSign := r.Group("/users")
	{
		userSign.POST("/signup", controller.SignUp)
		userSign.POST("/sign_in", controller.SignIn)
		userSign.POST("/forgot", controller.ForgetPassword)
		userSign.PATCH("/verify", controller.VerificationForPasswordChange)
	}

	users := r.Group("/users")
	users.Use(middleware.UserBasicAuth())
	{
		users.POST("/:user_id/details", controller.UserDetails)
	}
	payments := r.Group("/payments")
	payments.Use(middleware.UserBasicAuth())
	{
		payments.POST(":user_id/add_amount", controller.AddAmount)
		payments.POST(":user_id/withdraw_amount", controller.WithdrawAmount)
		payments.GET(":payment_status", controller.Callback)
	}

	reports := r.Group("/reports")
	reports.Use(middleware.UserBasicAuth())
	{
		reports.GET("pending_orders/:user_id", controller.DailyPendingOrders)
		reports.GET("holdings/:user_id", controller.Portfolio)
		reports.GET("order_history/:user_id", controller.OrdersHistory)
		reports.GET("profit_loss_history/:user_id", controller.ProfitLossHistory)
	}
	return r
}

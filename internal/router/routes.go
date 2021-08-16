package router

import (
	"Tradeasy/internal/controller"
	"Tradeasy/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	trade := r.Group("/pending_orders")
	trade.Use(middleware.UserVerificationAuth())
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
	websocket.Use(middleware.ExchangeBasicAuth())
	{
		websocket.GET("/stocks", controller.StockHandler)
		websocket.GET("/orders", controller.OrderHandler)
	}

	watchlist := r.Group("/user_watchlist")
	watchlist.Use(middleware.UserVerificationAuth())
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
	payments.GET(":payment_status", controller.Callback)
	payments.Use(middleware.UserVerificationAuth())
	{
		payments.POST(":user_id/add_amount", controller.AddAmount)
		payments.POST(":user_id/withdraw_amount", controller.WithdrawAmount)
	}

	reports := r.Group("/reports")
	reports.Use(middleware.UserVerificationAuth())
	{
		reports.GET("pending_orders/:user_id", controller.DailyPendingOrders)
		reports.GET("holdings/:user_id", controller.Portfolio)
		reports.GET("order_history/:user_id", controller.OrdersHistory)
		reports.GET("profit_loss_history/:user_id", controller.ProfitLossHistory)
	}
	return r
}

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
	watchlist.Use(middleware.UserBasicAuth())
	{
		watchlist.POST("", controller.CreateWatchlist)
		watchlist.POST("/:watchlist_id/add", controller.AddStockEntry)
		watchlist.DELETE("/:watchlist_id", controller.DeleteStockEntry)
		watchlist.PATCH("/sort", controller.SortWatchlist)
	}

	userSignUp := r.Group("/users")
	{
		userSignUp.POST("/signup", controller.SignUp)
	}

	users := r.Group("/users")
	users.Use(middleware.UserBasicAuth())
	{
		users.POST("/:user_id/details", controller.UserDetails)
		users.POST("/sign_in", controller.SignIn)
		users.POST("/forgot", controller.ForgetPassword)
		users.PATCH("/verify", controller.VerificationForPasswordChange)
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

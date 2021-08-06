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
	//users:=r.Group("/users")
	//{
	//	users.POST("/signup", controller.SignUp)
	//	users.POST("/:Userid/details", controller.UserDetails)
	//	users.POST("/signIn", controller.SignIn)
	//	users.POST("/forgot", controller.ForgetPassword)
	//	users.PATCH("/verify", controller.VerificationForPasswordChange)
	//}
	return r
}

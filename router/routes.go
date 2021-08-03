package router

import "pkg/mod/github.com/gin-gonic/gin@v1.7.2"

package routes

import (
Controllers "gin-framework"
"pkg/mod/github.com/gin-gonic/gin@v1.7.2"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	grp := r.Group("/users")
	{
		grp.POST("/signup", controller.SignUp)
		grp.POST("/:Userid/details", Controllers.UserDetails)
		grp.POST("/signIn", Controllers.SignIn)
		grp.POST("/forgot", Controllers.ForgetPassword)
		grp.PATCH("/verify", Controllers.VerificationForPasswordChange)
	}
	return r
}



package routes

import (
Controllers "gin-framework"
"pkg/mod/github.com/gin-gonic/gin@v1.7.2"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	grp := r.Group("/user_watchlist")
	{
		grp.POST("", Controllers.CreateWatchlist)
		grp.POST("/:watchlist_id/add", Controllers.AddStockEntry)
		grp.DELETE("/:watchlist_id", Controllers.DeleteStockEntry)
		grp.PATCH("/sort", Controllers.SortWatchlist)
	}
	return r
}





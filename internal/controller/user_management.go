package controller

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/user_management"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	var user model.Users
	er:=c.BindJSON(&user)
	if er !=nil {
		return
	}
	var req user_management.SignUpRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res,err := user_management.SignUp(&user,req)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func UserDetails(c *gin.Context) {
	var user model.TradingAccount
	er:=c.BindJSON(&user)
	if er !=nil {
		return
	}
	var req user_management.UserDetailsRequest
	userid := c.Params.ByName("Userid")
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res,errUserDetails := user_management.UserDetails(&user,req,userid)
	if errUserDetails != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func SignIn(c *gin.Context) {
	var user model.Users
	er:=c.BindJSON(&user)
	if er !=nil {
		return
	}
	var req user_management.SignInRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res,errUserSignIn := user_management.UserSignIn(&user,req)
	if errUserSignIn != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}
func ForgetPassword(c *gin.Context) {
	var user model.Users
	er:=c.BindJSON(&user)
	if er !=nil {
		return
	}
	var req user_management.ForgetPasswordRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res,errForgetPassword := user_management.ForgetPassword(&user,req)
	if errForgetPassword != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func VerificationForPasswordChange(c *gin.Context) {
	var req user_management.VerifyRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res,errVer := user_management.VerificationForPasswordChange(req)
	if errVer != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}


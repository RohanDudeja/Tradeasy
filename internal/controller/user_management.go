package controller

import (
	"Tradeasy/internal/services/user_management"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SignUp(c *gin.Context) {
	var req user_management.SignUpRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, err := user_management.SignUp(req)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func UserDetails(c *gin.Context) {
	var req user_management.UserDetailsRequest
	userid := c.Params.ByName("user_id")
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, err := user_management.UserDetails(req, userid)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func SignIn(c *gin.Context) {
	var req user_management.SignInRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, err := user_management.UserSignIn(req)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}
func ForgetPassword(c *gin.Context) {
	var req user_management.ForgetPasswordRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, err := user_management.ForgetPassword(req)
	if err != nil {
		log.Println(err.Error())
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
	res, err := user_management.VerificationForPasswordChange(req)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

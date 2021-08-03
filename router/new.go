package router

import (
	"Tradeasy/model"
	"fmt"
	"gin-microservice/Models"
	"math/big"
	"net/http"
	"pkg/mod/github.com/gin-gonic/gin@v1.7.2"
	"strings"
	"time"
)

func SignUp(c *gin.Context) {
	var user model.Users
	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	email := user.EmailId
	user.Userid = strings.Split(email, "@")[0]
	user.CreatedAt = time.Now()
	errSignUp := Users.SignUp(&user)
	if errSignUp != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func UserDetails(c *gin.Context) {
	var user model.UsersTradingAccDetails
	userid := c.Params.ByName("Userid")
	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	user.Userid = userid
	user.Balance = 0
	user.CreatedAt = time.Now()
	errUserDetails := Users.UserDetails(&user)
	if errUserDetails != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func UserSignIn(c *gin.Context) {
	var user Models.Users
	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	errUserSignIn := Users.UserSignIn(&user)
	if errUserSignIn != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
func ForgetPassword(c *gin.Context) {
	var user Models.Users
	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	errForgetPassword := Users.GetUserByUserid(&user)
	if errForgetPassword != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func VerificationForPasswordChange(c *gin.Context) {
	var user Models.Users
	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	errVerificationForPasswordChange := Users.VerificationForPasswordChange(&user)
	if errVerificationForPasswordChange != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

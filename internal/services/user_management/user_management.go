package user_management

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/provider/redis"
	"Tradeasy/internal/utils"
	"errors"
	"strings"
	"time"
)

var OTPExpiryTime = 5 * time.Minute

func SignUp(req SignUpRequest) (res SignUpResponse, err error) {
	email := req.EmailId

	var user model.Users
	err = database.GetDB().Table("users").Where("email_id = ?", req.EmailId).First(&user).Error
	if err == nil {
		return res, errors.New("email id already registered")
	}
	err = database.GetDB().Table("users").Where("password = ?", res.Password).First(&user).Error
	if err == nil {
		return res, errors.New("password already taken")
	}
	user.UserId = strings.Split(email, "@")[0]
	user.EmailId = req.EmailId
	user.Password = req.Password
	err = database.GetDB().Table("users").Create(&user).Error
	if err != nil {
		return res, errors.New("signUp failed")
	}
	res.UserId = user.UserId
	res.Password = user.Password
	res.Message = "User registered"
	return res, nil
}

func UserDetails(req UserDetailsRequest, userid string) (res UserDetailsResponse, err error) {

	var user model.Users
	var tradingAcc model.TradingAccount

	err = database.GetDB().Table("users").Where("user_id = ? ", userid).First(&user).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	err = database.GetDB().Table("trading_account").Where("user_id = ? OR pan_card_no = ? OR bank_acc_no = ?", userid, req.PanCardNo, req.BankAccNo).First(&tradingAcc).Error
	if err == nil {
		return res, errors.New("user details already registered")
	}
	tradingAcc.UserId = userid
	tradingAcc.PanCardNo = req.PanCardNo
	tradingAcc.BankAccNo = req.BankAccNo
	tradingAcc.TradingAccId = "TRA" + userid
	tradingAcc.Balance = 0

	err = database.GetDB().Table("trading_account").Create(&tradingAcc).Error
	if err != nil {
		return res, errors.New("user details failed to enter")
	}
	res.TradingAccId = tradingAcc.TradingAccId
	res.Balance = tradingAcc.Balance
	res.Message = "User Details registered"
	return res, nil
}
func UserSignIn(req SignInRequest) (res SignInResponse, err error) {

	var user model.Users
	err = database.GetDB().Table("users").Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	err = database.GetDB().Table("users").Where("user_id = ? AND password = ?", req.UserId, req.Password).First(&user).Error
	if err != nil {
		return res, errors.New("incorrect password")
	}
	res.Message = "Signed in successfully"
	return res, nil
}

func ForgetPassword(req ForgetPasswordRequest) (res ForgetPasswordResponse, err error) {

	var user model.Users
	err = database.GetDB().Table("users").Where("user_id = ? AND email_id = ?", req.UserId, req.EmailId).First(&user).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	otp, err := utils.GetRandNum()
	if err != nil {
		return res, errors.New("otp not generated")
	}
	err = redis.SetValue(req.EmailId, otp, OTPExpiryTime)
	if err != nil {
		return res, errors.New("otp not generated")
	}
	res.Otp = otp
	return res, nil
}

func VerificationForPasswordChange(req VerifyRequest) (res VerifyResponse, err error) {
	originalOtp, err := redis.GetValue(req.EmailId)
	if err != nil {
		return res, errors.New("verification failed")
	}
	if req.Otp != originalOtp {
		return res, errors.New("verification failed")
	}
	database.GetDB().Table("users").Where("user_id = ? AND email_id = ?", req.UserId, req.EmailId).Update("password", req.NewPassword)
	res.UserId = req.UserId
	res.NewPassword = req.NewPassword
	res.Message = "Password changed successfully"
	return res, nil
}

func CheckUserDetails(userid string) error {
	err := database.GetDB().Table("trading_account").Where("user_id = ? ", userid).Error
	if err != nil {
		return errors.New("add your details first")
	}
	return nil
}

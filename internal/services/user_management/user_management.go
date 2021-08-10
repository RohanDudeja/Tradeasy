package user_management

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/redis"
	"Tradeasy/internal/utils"
	"errors"
	"strings"
	"time"
)

var OTPExpiryTime = 5 * time.Minute

func SignUp(Req SignUpRequest) (Res SignUpResponse, err error) {
	email := Req.EmailId

	var user model.Users
	err = config.DB.Table("users").Where("email_id = ?", Req.EmailId).First(&user).Error
	if err == nil {
		return Res, errors.New("email id already registered")
	}
	err = config.DB.Table("users").Where("password = ?", Res.Password).First(&user).Error
	if err == nil {
		return Res, errors.New("password already taken")
	}
	user.UserId = strings.Split(email, "@")[0]
	user.EmailId = Req.EmailId
	user.Password = Req.Password
	err = config.DB.Table("users").Create(&user).Error
	if err != nil {
		return Res, errors.New("signUp failed")
	}
	Res.UserId = strings.Split(email, "@")[0]
	Res.Password = Req.Password
	Res.Message = "User registered"
	return Res, nil
}

func UserDetails(Req UserDetailsRequest, userid string) (Res UserDetailsResponse, err error) {

	var user model.Users
	var tradingAcc model.TradingAccount

	err = config.DB.Table("users").Where("user_id = ? ", userid).First(&user).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	err = config.DB.Table("trading_account").Where("user_id = ? OR pan_card_no = ? OR bank_acc_no = ?", userid, Req.PanCardNo, Req.BankAccNo).First(&tradingAcc).Error
	if err == nil {
		return Res, errors.New("user details already registered")
	}
	tradingAcc.UserId = userid
	tradingAcc.PanCardNo = Req.PanCardNo
	tradingAcc.BankAccNo = Req.BankAccNo
	tradingAcc.TradingAccId = "TRA" + userid
	tradingAcc.Balance = 0

	err = config.DB.Table("trading_account").Create(&tradingAcc).Error
	if err != nil {
		return Res, errors.New("user details failed to enter")
	}
	Res.TradingAccId = "TRA" + userid
	Res.Balance = 0
	Res.Message = "User Details registered"
	return Res, nil
}
func UserSignIn(Req SignInRequest) (Res SignInResponse, err error) {

	var user model.Users
	err = config.DB.Table("users").Where("user_id = ?", Req.UserId).First(&user).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	err = config.DB.Table("users").Where("user_id = ? AND password = ?", Req.UserId, Req.Password).First(&user).Error
	if err != nil {
		return Res, errors.New("incorrect password")
	}
	Res.Message = "Signed in successfully"
	return Res, nil
}

func ForgetPassword(Req ForgetPasswordRequest) (Res ForgetPasswordResponse, err error) {

	var user model.Users
	err = config.DB.Table("users").Where("user_id = ? AND email_id = ?", Req.UserId, Req.EmailId).First(&user).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	otp, err := utils.GetRandNum()
	if err != nil {
		return Res, errors.New("otp not generated")
	}
	err = redis.SetValue(Req.EmailId, OTPExpiryTime)
	if err != nil {
		return Res, errors.New("otp not generated")
	}
	Res.Otp = otp
	return Res, nil
}

func VerificationForPasswordChange(Req VerifyRequest) (Res VerifyResponse, err error) {
	originalOtp, err := redis.GetValue(Req.EmailId)
	if err != nil {
		return Res, errors.New("verification failed")
	}
	if Req.Otp != originalOtp {
		return Res, errors.New("verification failed")
	}
	config.DB.Table("users").Where("user_id = ? AND email_id = ?", Req.UserId, Req.EmailId).Update("password", Res.NewPassword)
	Res.UserId = Req.UserId
	Res.NewPassword = Req.NewPassword
	Res.Message = "Password changed successfully"
	return Res, nil
}

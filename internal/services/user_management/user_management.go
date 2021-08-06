package user_management

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/redis"
	"errors"
	"strings"
	"time"
)

func SignUp(SUpReq SignUpRequest) (SUpRes SignUpResponse, err error) {
	email := SUpReq.EmailId
	SUpRes.UserId = strings.Split(email, "@")[0]
	SUpRes.Password = SUpReq.Password
	SUpRes.Message = "User registered"

	var user model.Users
	user.UserId = SUpRes.UserId
	user.EmailId = SUpReq.EmailId
	user.Password = SUpReq.Password
	err = config.DB.Table("users").Create(&user).Error
	if err != nil {
		return SUpRes, errors.New("signUp failed")
	}
	return SUpRes, nil
}

func UserDetails(detReq UserDetailsRequest, userid string) (detRes UserDetailsResponse, err error) {
	detRes.TradingAccId = "TRA" + userid
	detRes.Balance = 0
	detRes.Message = "User Details registered"

	var user model.TradingAccount
	user.UserId = userid
	user.PanCardNo = detReq.PanCardNo
	user.BankAccNo = detReq.BankAccNo
	user.TradingAccId = detRes.TradingAccId
	user.Balance = detRes.Balance

	err = config.DB.Table("users_trading_acc_details").Create(&user).Error
	if err != nil {
		return detRes, errors.New("enter correct user details")
	}
	return detRes, nil
}
func UserSignIn(SInReq SignInRequest) (SInRes SignInResponse, err error) {
	SInRes.Message = "Signed in successfully"

	var user model.Users
	err = config.DB.Table("users").Where("userid = ? AND password = ?", SInReq.UserId, SInReq.Password).First(&user).Error
	if err != nil {
		return SInRes, errors.New("sign in Failed")
	}
	return SInRes, nil
}

func ForgetPassword(FPReq ForgetPasswordRequest) (FPRes ForgetPasswordResponse, err error) {

	var user model.Users
	err = config.DB.Table("users").Where("userid = ? AND emailId = ?", FPReq.UserId, FPReq.EmailId).First(&user).Error
	if err != nil {
		return FPRes, errors.New("user not found")
	}
	otp, err_ := redis.GetRandNum()
	if err_ != nil {
		return FPRes, errors.New("otp not generated")
	}
	FPRes.Otp = otp
	e := redis.SetValue(FPReq.EmailId, otp, 5*time.Minute)
	if e != nil {
		return FPRes, errors.New("otp not generated")
	}
	return FPRes, nil
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse, err error) {
	VerRes.UserId = VerReq.UserId
	VerRes.NewPassword = VerReq.NewPassword
	VerRes.Message = "Password changed successfully"

	originalOtp, e := redis.GetValue(VerReq.EmailId)
	if e != nil {
		return VerRes, errors.New("verification failed")
	}
	if VerReq.Otp != originalOtp {
		return VerRes, errors.New("verification failed")
	}
	config.DB.Table("users").Where("userid = ? AND emailId = ?", VerReq.UserId, VerReq.EmailId).Update("password", VerRes.NewPassword)

	return VerRes, nil
}

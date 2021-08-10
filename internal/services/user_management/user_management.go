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

var ExpiryTime = 5 * time.Minute

func SignUp(SUpReq SignUpRequest) (SUpRes SignUpResponse, err error) {
	email := SUpReq.EmailId

	var user model.Users
	err = config.DB.Table("users").Where("email_id = ?", SUpReq.EmailId).First(&user).Error
	if err == nil {
		return SUpRes, errors.New("email id already registered")
	}
	err = config.DB.Table("users").Where("password = ?", SUpRes.Password).First(&user).Error
	if err == nil {
		return SUpRes, errors.New("password already taken")
	}
	user.UserId = SUpRes.UserId
	user.EmailId = SUpReq.EmailId
	user.Password = SUpReq.Password
	err = config.DB.Table("users").Create(&user).Error
	if err != nil {
		return SUpRes, errors.New("signUp failed")
	}
	SUpRes.UserId = strings.Split(email, "@")[0]
	SUpRes.Password = SUpReq.Password
	SUpRes.Message = "User registered"
	return SUpRes, nil
}

func UserDetails(detReq UserDetailsRequest, userid string) (detRes UserDetailsResponse, err error) {

	var user model.Users
	var ta model.TradingAccount

	err = config.DB.Table("users").Where("user_id = ? ", userid).First(&user).Error
	if err != nil {
		return detRes, errors.New("user not found")
	}
	err = config.DB.Table("trading_account").Where("user_id = ? OR pan_card_no = ? OR bank_acc_no = ?", userid, detReq.PanCardNo, detReq.BankAccNo).First(&ta).Error
	if err == nil {
		return detRes, errors.New("user details already registered")
	}
	ta.UserId = userid
	ta.PanCardNo = detReq.PanCardNo
	ta.BankAccNo = detReq.BankAccNo
	ta.TradingAccId = detRes.TradingAccId
	ta.Balance = detRes.Balance

	err = config.DB.Table("trading_account").Create(&ta).Error
	if err != nil {
		return detRes, errors.New("user details failed to enter")
	}
	detRes.TradingAccId = "TRA" + userid
	detRes.Balance = 0
	detRes.Message = "User Details registered"
	return detRes, nil
}
func UserSignIn(SInReq SignInRequest) (SInRes SignInResponse, err error) {

	var user model.Users
	err = config.DB.Table("users").Where("user_id = ?", SInReq.UserId).First(&user).Error
	if err != nil {
		return SInRes, errors.New("user not found")
	}
	err = config.DB.Table("users").Where("password = ?", SInReq.Password).First(&user).Error
	if err != nil {
		return SInRes, errors.New("incorrect password")
	}
	SInRes.Message = "Signed in successfully"
	return SInRes, nil
}

func ForgetPassword(FPReq ForgetPasswordRequest) (FPRes ForgetPasswordResponse, err error) {

	var user model.Users
	err = config.DB.Table("users").Where("user_id = ? AND email_id = ?", FPReq.UserId, FPReq.EmailId).First(&user).Error
	if err != nil {
		return FPRes, errors.New("user not found")
	}
	otp, err_ := utils.GetRandNum()
	if err_ != nil {
		return FPRes, errors.New("otp not generated")
	}
	err = redis.SetValue(FPReq.EmailId, ExpiryTime)
	if err != nil {
		return FPRes, errors.New("otp not generated")
	}
	FPRes.Otp = otp
	return FPRes, nil
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse, err error) {
	originalOtp, err := redis.GetValue(VerReq.EmailId)
	if err != nil {
		return VerRes, errors.New("verification failed")
	}
	if VerReq.Otp != originalOtp {
		return VerRes, errors.New("verification failed")
	}
	config.DB.Table("users").Where("user_id = ? AND email_id = ?", VerReq.UserId, VerReq.EmailId).Update("password", VerRes.NewPassword)
	VerRes.UserId = VerReq.UserId
	VerRes.NewPassword = VerReq.NewPassword
	VerRes.Message = "Password changed successfully"
	return VerRes, nil
}

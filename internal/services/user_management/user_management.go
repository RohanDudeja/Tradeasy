package user_management

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"strings"
)

func SignUp(user *model.Users,SUpReq SignUpRequest) (SUpRes SignUpResponse,err error) {
	email := SUpReq.EmailId
	SUpRes.UserId = strings.Split(email, "@")[0]
	SUpRes.Password=SUpReq.Password
	SUpRes.Message="User registered"

	user.UserId=SUpRes.UserId

	err = config.DB.Table("users").Create(user).Error
	return SUpRes,err
}

func GetUserByUserid(user *model.TradingAccount,userid string) (err error) {
	if err = config.DB.Table("users").Where("userid = ?", userid).First(&user).Error; err != nil {
		return err
	}
	return nil
}
func UserDetails(user *model.TradingAccount,detReq UserDetailsRequest,userid string) (detRes UserDetailsResponse,err error) {
	detRes.TradingAccId="TRA"+userid
	detRes.Balance = 0
	detRes.Message ="User Details registered"

	user.UserId=userid
	user.TradingAccId=detRes.TradingAccId
	user.Balance=detRes.Balance

	err = config.DB.Table("users_trading_acc_details").Create(user).Error

	return detRes,err
}
func UserSignIn(user* model.Users,SInReq SignInRequest) (SInRes SignInResponse,err error) {
	SInRes.Message="Signed in successfully"

	err = config.DB.Table("users").Where("userid = ? AND password = ?", user.UserId, user.Password).First(user).Error

	return SInRes,err
}

func ForgetPassword(user *model.Users,FPReq ForgetPasswordRequest) (FPRes ForgetPasswordResponse,err error) {

	err = config.DB.Table("users").Where("userid = ? AND emailId = ?", user.UserId, user.EmailId).First(user).Error
	return FPRes,err
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse,err error) {
	return VerRes,err
}

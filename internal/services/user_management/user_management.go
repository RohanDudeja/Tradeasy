package user_management

import (
	"strings"
)

func SignUp(SUpReq SignUpRequest) (SUpRes SignUpResponse,err error) {
	email := SUpReq.EmailId
	SUpRes.UserId = strings.Split(email, "@")[0]
	SUpRes.Password=SUpReq.Password
	SUpRes.Message="User registered"
	//user.CreatedAt = time.Now()
	//err = Config.DB.Table("users").Create(user).Error
	//if err != nil {
	//	return err
	//}
	return SUpRes,err
}

// UserDetails func GetUserByUserid(user *model.Users,userid string) (err error) {
//	if err = config.DB.Table("users").Where("userid = ?", userid).First(&user).Error; err != nil {
//		return err
//	}
//	return nil
//}
func UserDetails(detReq UserDetailsRequest,userid string) (detRes UserDetailsResponse,err error) {
	detRes.TradingAccId="TRA"+userid
	//detRes.Balance = 0
	detRes.Message ="User Details registered"
	//user.CreatedAt = time.Now()
	//if err = Config.DB.Table("users_trading_acc_details").Create(user).Error; err != nil {
	//	return err
	//}
	return detRes,err
}
func UserSignIn(SInReq SignInRequest) (SInRes SignInResponse,err error) {
	SInRes.Message="Signed in successfully"
	//if err = Config.DB.Table("users").Where("userid = ? AND password = ?", user.Userid, user.Password).First(user).Error; err != nil {
	//	return err
	//}
	return SInRes,err
}

func ForgetPassword(FPReq ForgetPasswordRequest) (FPRes ForgetPasswordResponse,err error) {
	return FPRes,err
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse,err error) {
	return VerRes,err
}

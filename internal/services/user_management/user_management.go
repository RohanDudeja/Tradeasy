package user_management

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"crypto/rand"
	"errors"
	"github.com/go-redis/redis"
	"github.com/nu7hatch/gouuid"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:3306",
	Password: "password", // no password set
	DB:       0,  // use default DB
})

// SetValue sets the key value pair
func SetValue(key string, value string, expiry time.Duration) error {
	err := redisClient.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}
func getRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}
func randToken() (string, error) {
	// Using UUID V5 for generating the Token
	u4, err := uuid.NewV4()
	UUIDtoken := u4.String()
	if err != nil {
		return "", err
	}
	return UUIDtoken, nil
}
// GetValue the value corresponding to a given key
func GetValue(key string) (string, error) {
	value, argh := redisClient.Get(key).Result()
	if argh != nil {
		return "", argh
	}
	return value, nil
}


func SignUp(user *model.Users,SUpReq SignUpRequest) (SUpRes SignUpResponse,err error) {
	email := SUpReq.EmailId
	SUpRes.UserId = strings.Split(email, "@")[0]
	SUpRes.Password=SUpReq.Password
	SUpRes.Message="User registered"

	user.UserId=SUpRes.UserId

	err = config.DB.Table("users").Create(user).Error
	return SUpRes,err
}

// UserDetails func GetUserByUserid(user *model.TradingAccount,userid string) (err error) {
//	if err = config.DB.Table("users").Where("userid = ?", userid).First(&user).Error; err != nil {
//		return err
//	}
//	return nil
//}
func UserDetails(user *model.TradingAccount,detReq UserDetailsRequest,userid string) (detRes UserDetailsResponse,err error) {
	detRes.TradingAccId="TRA"+userid
	//detRes.Balance = big.NewInt(0)
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

	otp,err_:=getRandNum()
	token,er:=randToken()
	FPRes.Otp=otp
	FPRes.Nonce=token

	if err_ !=nil {
		return FPRes,err_
	}
	if er !=nil {
		return FPRes,er
	}
	e:= SetValue(FPReq.EmailId, otp, 5*time.Minute)
	e1:=SetValue(token,FPReq.EmailId,5*time.Minute)
	if e != nil{
		return FPRes,e
	}
	if e1 !=nil {
		return FPRes,e1
	}
	return FPRes,err
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse,err error) {

	VerRes.UserId=VerReq.UserId
	VerRes.Password=VerReq.Password

	emailId, e := GetValue(VerReq.Nonce)
	if e != nil {
		VerRes.Message="Verification failed"
		return VerRes,e
	}
	originalOtp, e := GetValue(emailId)
	if e != nil {
		VerRes.Message="Verification failed"
		return VerRes,e
	}
	if VerReq.Otp!=originalOtp {
		return VerRes,errors.New("verification failed")
	}
	config.DB.Table("users").Where("userid = ?",VerReq.UserId).Update("password",VerRes.Password)

	return VerRes,nil
}

package user_management

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"crypto/rand"
	"errors"
	"github.com/go-redis/redis"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
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

// GetValue the value corresponding to a given key
func GetValue(key string) (string, error) {
	value, argh := redisClient.Get(key).Result()
	if argh != nil {
		return "", argh
	}
	return value, nil
}

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
	otp, err_ := getRandNum()
	if err_ != nil {
		return FPRes, errors.New("otp not generated")
	}
	FPRes.Otp = otp
	e := SetValue(FPReq.EmailId, otp, 5*time.Minute)
	if e != nil {
		return FPRes, errors.New("otp not generated")
	}
	return FPRes, nil
}

func VerificationForPasswordChange(VerReq VerifyRequest) (VerRes VerifyResponse, err error) {
	VerRes.UserId = VerReq.UserId
	VerRes.NewPassword = VerReq.NewPassword
	VerRes.Message = "Password changed successfully"

	originalOtp, e := GetValue(VerReq.EmailId)
	if e != nil {
		return VerRes, errors.New("verification failed")
	}
	if VerReq.Otp != originalOtp {
		return VerRes, errors.New("verification failed")
	}
	config.DB.Table("users").Where("userid = ? AND emailId = ?", VerReq.UserId, VerReq.EmailId).Update("password", VerRes.NewPassword)

	return VerRes, nil
}

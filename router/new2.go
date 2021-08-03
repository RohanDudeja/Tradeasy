package router

func SignUp(user *Users) (err error) {
	if err = Config.DB.Table("users").Create(user).Error; err != nil {
		return err
	}
	return nil
}

func UserDetails(user *UsersTradingAccDetails) (err error) {
	if err = Config.DB.Table("users_trading_acc_details").Create(user).Error; err != nil {
		return err
	}
	return nil
}
func GetUserByUserid(user *Users) (err error) {
	if err = Config.DB.Table("users").Where("userid = ?", user.Userid).First(&user).Error; err != nil {
		return err
	}
	return nil
}
func UserSignIn(user *Users) (err error) {
	if err = Config.DB.Table("users").Where("userid = ? AND password = ?", user.Userid, user.Password).First(&user).Error; err != nil {
		return err
	}
	return nil
}

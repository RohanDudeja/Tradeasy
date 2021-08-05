package cmd

import (
	"Tradeasy/config"
	"Tradeasy/internal/router"
	"fmt"
	"github.com/jinzhu/gorm"
)
var err_ error
func main()  {
	config.DB,err_=gorm.Open("mysql",config.DbURL(config.BuildConfig()))
	if err_ != nil {
		fmt.Println("Status:", err_)
	}
	defer config.DB.Close()
	r:=router.SetUpRouter()
	r.Run()
}

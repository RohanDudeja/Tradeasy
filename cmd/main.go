package cmd

import "Tradeasy/internal/router"
func main()  {
	r:=router.SetUpRouter()
	r.Run()
}

package Order

func BuyOrder(breq BuyRequest) (err error, bres BuyResponse) {
	return err, bres
}

func SellOrder(sreq SellRequest) (err error, sres SellResponse) {
	return err, sres
}

func CancelOrder(id string) (err error, cres CancelResponse) {
	return err, cres
}

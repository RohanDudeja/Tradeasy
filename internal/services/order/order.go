package order

func BuyOrder(BReq BuyRequest) (BRes BuyResponse, err error) {
	return BRes, err
}

func SellOrder(BReq SellRequest) (SRes SellResponse, err error) {
	return SRes, err
}

func CancelOrder(id string) (CRes CancelResponse, err error) {
	return CRes, err
}

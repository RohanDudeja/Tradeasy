package watchlist

func CreateWatchlist(CrReq CreateRequest) (CrRes CreateResponse, err error) {
	return CrRes, err
}

func AddStockEntry(AddReq AddStockRequest, watchlistId string) (AddRes AddStockResponse, err error) {
	return AddRes, err
}

func DeleteStockEntry(DelReq DeleteStockRequest, watchlistId string) (DelRes DeleteStockResponse, err error) {
	return DelRes, err
}
func SortWatchlist(SortReq SortRequest) (SortRes SortResponse, err error) {
	return SortRes, err
}

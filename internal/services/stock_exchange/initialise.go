package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Results []Values `json:"results"`
}
type Values struct {
	Ticker string `json:"ticker"`
}
type StockFeed struct {
	Open      float64 `json:"open"`
	PrevClose float64 `json:"close"`
	Ticker    string  `json:"symbol"`
	LTP       float64 `json:"afterHours"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
}

func RandomizerAlgo() {

	var allStocks []model.Stocks
	config.DB.Model(&allStocks)
	orderType := []string{"Limit", "Market"}
	for _, stock := range allStocks {

		//placing buy order
		u := uuid.New().String()
		rand.Seed(time.Now().UnixNano())
		idx := rand.Intn(2)
		order := orderType[idx]
		min := stock.LTP - int(float64(stock.LTP)*0.01)
		max := stock.LTP + int(float64(stock.LTP)*0.01)
		buyOrderBody := OrderRequest{
			OrderID:         u,
			StockName:       stock.StockName,
			OrderPlacedTime: time.Time{},
			OrderType:       order,
			LimitPrice:      uint(rand.Intn(max-min+1) + min),
			Quantity:        uint(rand.Intn(100) + 1),
		}
		BuyOrder(buyOrderBody)

		//placing sell order
		u = uuid.New().String()
		rand.Seed(time.Now().UnixNano())
		idx = rand.Intn(2)
		order = orderType[idx]
		min = stock.LTP - int(float64(stock.LTP)*0.01)
		max = stock.LTP + int(float64(stock.LTP)*0.01)
		time.Sleep(1 * time.Second)
		sellOrderBody := OrderRequest{
			OrderID:         u,
			StockName:       stock.StockName,
			OrderPlacedTime: time.Time{},
			OrderType:       order,
			LimitPrice:      uint(rand.Intn(max-min+1) + min),
			Quantity:        uint(rand.Intn(100) + 1),
		}
		SellOrder(sellOrderBody)
	}
	// sleep and run again
	time.Sleep(300 * time.Second)
	RandomizerAlgo()
}

func GetTickers(limit int) (tickers []string) {
	urlTickers := "https://api.polygon.io/v3/reference/tickers?active=true&sort=primary_exchange&order=asc&limit=" + strconv.Itoa(limit) + "&apiKey=721mkXq0CBNvCMi5iyJ9E1gBRDiFcT8b"
	req, _ := http.NewRequest("GET", urlTickers, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	rawResponse := bytes.NewReader(body)
	decoder := json.NewDecoder(rawResponse)
	parsedResponse := Response{}
	decoder.Decode(&parsedResponse)

	for _, r := range parsedResponse.Results {
		//fmt.Println(r.Ticker)
		tickers = append(tickers, r.Ticker)
	}
	return tickers
}

func InitialiseStock(ticker string) {

	url1 := "https://api.polygon.io/v1/open-close/" + ticker + "/2021-08-04?adjusted=true&apiKey=721mkXq0CBNvCMi5iyJ9E1gBRDiFcT8b"
	req, _ := http.NewRequest("GET", url1, nil)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	rawResponse := bytes.NewReader(body)
	decoder := json.NewDecoder(rawResponse)
	parsedResponse := StockFeed{}
	decoder.Decode(&parsedResponse)
	newStock := model.Stocks{
		StockTickerSymbol: parsedResponse.Ticker,
		StockName:         parsedResponse.Ticker,
		LTP:               int(parsedResponse.LTP * 100.0),
		OpenPrice:         int(parsedResponse.Open * 100.0),
		HighPrice:         int(parsedResponse.High * 100.0),
		LowPrice:          int(parsedResponse.Low * 100.0),
		PreviousDayClose:  int(parsedResponse.PrevClose * 100.0),
		PercentageChange:  int(100.0 * (parsedResponse.PrevClose - parsedResponse.LTP) / parsedResponse.PrevClose),
	}
	config.DB.Create(&newStock)
}

func InitialiseAllStocks() {
	stocksNeeded := 20
	tickers := GetTickers(stocksNeeded)
	for _, ticker := range tickers {
		InitialiseStock(ticker)
	}
}

func CreateBuyersAndSellers(ticker string, quantity int, ltp int) {

	rand.Seed(time.Now().UnixNano())
	min := ltp - int(float64(ltp)*0.01)
	max := ltp + int(float64(ltp)*0.01)
	newBuy := model.BuyOrderBook{
		OrderID:           uuid.New().String(),
		StockTickerSymbol: ticker,
		OrderQuantity:     quantity,
		OrderStatus:       "Pending",
		OrderPrice:        rand.Intn(max-min+1) + min,
	}
	config.DB.Create(&newBuy)

	time.Sleep(1 * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	min = ltp - int(float64(ltp)*0.01)
	max = ltp + int(float64(ltp)*0.01)
	newSell := model.SellOrderBook{
		OrderID:           uuid.New().String(),
		StockTickerSymbol: ticker,
		OrderQuantity:     quantity,
		OrderStatus:       "Pending",
		OrderPrice:        rand.Intn(max-min+1) + min,
	}
	config.DB.Create(&newSell)
}

func InitialiseBuyersAndSellers() {

	var allStocks []model.Stocks
	config.DB.Find(&allStocks)

	quantities := []int{100, 50, 120, 200, 280}
	for _, stock := range allStocks {
		for _, quantity := range quantities {
			CreateBuyersAndSellers(stock.StockTickerSymbol, quantity, stock.LTP)
		}
	}

}

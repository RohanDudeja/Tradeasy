package stock_exchange

import (
	model "Tradeasy/internal/model/stock_exchange"
	"Tradeasy/internal/provider/database"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

const (
	PercentChange       = 0.01
	OrdersQuantityRange = 100
	StocksNeeded        = 20
)

// GetTickers ...Get tickers as per requirement to be used stocks initialisation
func GetTickers(limit int) (tickers []string, err error) {

	apiKey := "721mkXq0CBNvCMi5iyJ9E1gBRDiFcT8b"
	baseURL := "https://api.polygon.io/v3/reference/tickers?market=stocks&sort=last_updated_utc&order=asc&limit="
	urlTickers := baseURL + strconv.Itoa(limit) + "&apiKey=" + apiKey
	//send a get request
	req, _ := http.NewRequest("GET", urlTickers, nil)
	res, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	rawResponse := bytes.NewReader(body)
	decoder := json.NewDecoder(rawResponse)
	parsedResponse := Response{}
	err = decoder.Decode(&parsedResponse)
	if err != nil {
		log.Println(err.Error())
		return tickers, err
	}

	for _, r := range parsedResponse.Results {
		tickers = append(tickers, r.Ticker)
	}
	return tickers, err
}

// InitialiseStock ...fetch data for a given stock in arguments
func InitialiseStock(ticker string) {

	var existingStocks []model.Stocks
	err := database.GetDB().Table("stocks").Select("*").Where("stock_name = ?", ticker).Find(&existingStocks).Error
	if err != nil {
		log.Println(err.Error())
	}
	if len(existingStocks) == 0 {
		baseURL := "https://api.polygon.io/v1/open-close/"
		date := "/2021-08-11"
		apiKey := "721mkXq0CBNvCMi5iyJ9E1gBRDiFcT8b"
		stocksURL := baseURL + ticker + date + "?adjusted=true&apiKey=" + apiKey
		req, _ := http.NewRequest("GET", stocksURL, nil)

		res, _ := http.DefaultClient.Do(req)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err.Error())
				return
			}
		}(res.Body)
		body, _ := ioutil.ReadAll(res.Body)
		rawResponse := bytes.NewReader(body)
		decoder := json.NewDecoder(rawResponse)
		parsedResponse := StockFeed{}
		err := decoder.Decode(&parsedResponse)
		if err != nil {
			log.Println(err.Error())
			return
		}
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
		if parsedResponse.Ticker != "" {
			err := database.GetDB().Create(&newStock).Error
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

// InitialiseAllStocks ...initialise all data points in the database for every stock
func InitialiseAllStocks() {
	tickers, err := GetTickers(StocksNeeded)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, ticker := range tickers {
		InitialiseStock(ticker)
	}
}

func InitialiseRandomizer() {
	go RandomizerAlgo()
}

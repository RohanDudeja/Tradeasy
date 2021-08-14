package stock_exchange

import (
	"Tradeasy/internal/provider/database"
	"Tradeasy/test/mysql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteBuyOrder_Success(t *testing.T) {
	sqlDB, mock := mysql.NewMock()
	db, err := gorm.Open("mysql", sqlDB)
	database.SetDB(db)
	assert.NoError(t, err)

	query := "DELETE FROM buy_order_book WHERE order_id = ?"

	mock.ExpectExec(query).WithArgs("ord123").WillReturnResult(sqlmock.NewResult(0, 1))

	delRes, err := DeleteBuyOrder("ord123")
	if assert.NoError(t, err) {
		assert.Equal(t, delRes.Success, true)
	} else if assert.Error(t, err) {
		assert.Equal(t, delRes.Success, false)
	}
}

//func TestDeleteSellOrder(t *testing.T) {
//	sqlDB, mock := mysql.NewMock()
//	db, err := gorm.Open("mysql", sqlDB)
//	database.SetDB(db)
//	assert.NoError(t, err)
//
//	query := "DELETE FROM sell_order_book WHERE order_id = ?"
//	mock.ExpectExec(query).WithArgs("ord123").WillReturnResult(sqlmock.NewResult(0, 1))
//
//	delRes, err := DeleteSellOrder("ord123")
//	if assert.NoError(t, err) {
//		assert.Equal(t, delRes.Success, true)
//	} else if assert.Error(t, err) {
//		assert.Equal(t, delRes.Success, false)
//	}
//}
//
//func TestGetLTP(t *testing.T) {
//	sqlDB, mock := mysql.NewMock()
//	db, err := gorm.Open("mysql", sqlDB)
//	database.SetDB(db)
//	assert.NoError(t, err)
//	query := "SELECT * FROM stocks WHERE stock_ticker_symbol = ?"
//	stocksTable := model.Stocks{
//		ID:                1,
//		StockTickerSymbol: "AAA",
//		StockName:         "Ab",
//		LTP:               100,
//		OpenPrice:         99,
//		HighPrice:         120,
//		LowPrice:          90,
//		PreviousDayClose:  80,
//		PercentageChange:  100,
//		CreatedAt:         time.Now(),
//		UpdatedAt:         time.Now(),
//	}
//	rows := sqlmock.NewRows([]string{"id", "stock_ticker_symbol", "stock_name", "ltp",
//		"open_price", "high_price", "low_price", "previous_day_close", "percentage_change",
//		"created_at", "updated_at"}).
//		AddRow(stocksTable.ID, stocksTable.StockTickerSymbol, stocksTable.StockName, stocksTable.LTP,
//			stocksTable.OpenPrice, stocksTable.HighPrice, stocksTable.LowPrice, stocksTable.PreviousDayClose,
//			stocksTable.PercentageChange, stocksTable.CreatedAt, stocksTable.UpdatedAt)
//	mock.ExpectQuery(query).WithArgs(stocksTable.StockTickerSymbol).WillReturnRows(rows)
//	value, err := GetLTP(stocksTable.StockTickerSymbol)
//	assert.NoError(t, err)
//	assert.Equal(t, value, stocksTable.LTP)
//}

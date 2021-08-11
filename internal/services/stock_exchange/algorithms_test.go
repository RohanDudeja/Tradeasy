package stock_exchange

import (
	"Tradeasy/internal/provider/database"
	"Tradeasy/test/mysql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteBuyOrder(t *testing.T) {
	mysql.NewMock()
	db, err := gorm.Open("mysql", mysql.GetSqlDB())
	database.SetDB(db)
	assert.NoError(t, err)
	mock := mysql.GetSqlMock()

	query := "DELETE FROM buy_order_book WHERE order_id = ?"

	mock.ExpectExec(query).WithArgs("ord123").WillReturnResult(sqlmock.NewResult(0, 1))

	_, err = DeleteBuyOrder("ord123")
	assert.NoError(t, err)
}

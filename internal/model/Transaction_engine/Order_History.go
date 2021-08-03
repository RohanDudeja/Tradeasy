package Transaction_engine

import "github.com/golang/protobuf/ptypes/timestamp"

type OrderHistory struct {
	Userid        string              `json:"userid"`
	User          Users               `gorm:"foreign_key:Userid"`
	OrderId       string              `gorm:"column:order_id" json:"order_id"`
	PendingOrders PendingOrders       `gorm:"foreign_key:OrderId"`
	Id            int                 `gorm:"primaryKey; column:id" json:"id"`
	StockName     string              `json:"stock_name" gorm:"column:stock_name""`
	Quantity      int                 `json:"quantity" gorm:"column:quantity"`
	BuyPrice      int                 `json:"buy_price" gorm:"column:buy_price"`
	SellPrice     int                 `json:"sell_price" gorm:"column:sell_price"`
	CommissionFee int                 `json:"commission_fee" gorm:"column:commission_fee"`
	BoughtAt      timestamp.Timestamp `json:"bought_at" gorm:"column:bought_at"`
	SoldAt        timestamp.Timestamp `json:"sold_at" gorm:"column:sold_at"`
	CreatedAt     timestamp.Timestamp `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     timestamp.Timestamp `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt timestamp.Timestamp `json:"deleted_at" gorm:"column:deleted_at"`
}


func (o *OrderHistory) TableName() string  {
	return "OrderHistory"
}



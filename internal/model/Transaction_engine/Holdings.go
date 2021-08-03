package Transaction_engine

import "github.com/golang/protobuf/ptypes/timestamp"

type Holdings struct{
	Userid        string
	User          Users               `gorm:"foreign_key:Userid" json:"user"`
	OrderId       string              `gorm:"column:order_id" json:"order_id"`
	PendingOrders PendingOrders       `gorm:"foreign_key:OrderId" json:"pending_orders"`
	Id            int                 `gorm:"primary_key; column:id" json:"id"`
	StockName     string              `gorm:"column:stock_name" json:"stock_name"`
	Quantity      int                 `json:"quantity" gorm:"column:quantity"`
	BuyPrice      int                 `json:"buy_price" gorm:"column:buy_price"`
	OrderedAt     timestamp.Timestamp `json:"ordered_at" gorm:"column:ordered_at"`
	CreatedAt     timestamp.Timestamp `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     timestamp.Timestamp `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt     timestamp.Timestamp `json:"deleted_at" gorm:"column:deleted_at"`
}

func (h *Holdings) TableName() string  {
	return "Holdings"
}

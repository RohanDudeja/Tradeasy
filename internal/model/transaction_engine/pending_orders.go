package transaction_engine

import "github.com/golang/protobuf/ptypes/timestamp"

type PendingOrders struct {
	Userid     string              `json:"userid" gorm:"foreign_key; column:userid"`
	OrderId    string              `gorm:"primary_key;column:order_id" json:"order_id"`
	StockName  string              `gorm:"column:stock_name" json:"stock_name"`
	OrderType  string              `gorm:"column:order_type" json:"order_type"`
	BookType   string              `gorm:"column:book_type" json:"book_type"`
	LimitPrice int                 `json:"limit_price" gorm:"column:limit_price"`
	Quantity   int                 `json:"quantity" gorm:"column:quantity"`
	OrderPrice int                 `json:"order_price" gorm:"column:order_price"`
	Status     string              `gorm:"column:status" json:"status"`
	CreatedAt  timestamp.Timestamp `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  timestamp.Timestamp `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  timestamp.Timestamp `json:"deleted_at" gorm:"column:deleted_at"`
}

func (p *PendingOrders) TableName() string {
	return "PendingOrders"
}

package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"time"
)

func UpdateBuyOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	var account model.TradingAccount

	if res.Status == status[3] {
		//Buy order Failed
		p.Status = status[3]
		if err = config.DB.Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}
		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			return err
		}
		account.Balance = account.Balance + int64(p.Quantity*p.OrderPrice)
		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			return err
		}
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Delete(p).Error; err != nil {
			return err
		}
		return nil
	} else if res.Status == status[1] {
		//Buy order completed
		p.Status = status[1]
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}
		h := model.Holdings{
			UserId:    p.UserId,
			OrderId:   p.OrderId,
			StockName: p.StockName,
			Quantity:  p.Quantity,
			BuyPrice:  res.AveragePrice,
			OrderedAt: p.CreatedAt,
			CreatedAt: time.Now(),
			UpdatedAt: res.OrderExecutionTime,
		}
		if p.BookType == "Market" {
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return err
			}
			account.Balance = account.Balance + int64((p.OrderPrice-res.AveragePrice)*res.Quantity)
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return err
			}
		} else {
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return err
			}
			account.Balance = account.Balance + int64((p.LimitPrice-res.AveragePrice)*res.Quantity)
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return err
			}
		}

		if err = config.DB.Table("holdings").Create(h).Error; err != nil {
			return err
		}
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Delete(p).Error; err != nil {
			return err
		}
	} else if res.Status == status[2] {
		//Buy order Half Completed
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).First(&p).Error; err != nil {
			return err
		}
		p.Status = status[2]
		p.Quantity = p.Quantity - res.Quantity
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}

		h := model.Holdings{
			UserId:    p.UserId,
			OrderId:   p.OrderId,
			StockName: p.StockName,
			Quantity:  res.Quantity,
			BuyPrice:  res.AveragePrice,
			OrderedAt: p.CreatedAt,
			CreatedAt: time.Now(),
			UpdatedAt: res.OrderExecutionTime,
		}
		if p.BookType == "Market" {
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return err
			}
			account.Balance = account.Balance + int64((p.OrderPrice-res.AveragePrice)*res.Quantity)
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return err
			}
		} else {
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return err
			}
			account.Balance = account.Balance + int64((p.LimitPrice-res.AveragePrice)*res.Quantity)
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return err
			}
		}
		if err = config.DB.Table("holdings").Create(h).Error; err != nil {
			return err
		}
	}
	return nil
}

func UpdateSellOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	var account model.TradingAccount

	if res.Status == status[3] {
		//Sell order failed
		p.Status = status[3]
		if err = config.DB.Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Delete(p).Error; err != nil {
			return err
		}
		return nil
	} else if res.Status == status[1] {
		//Sell order Completed
		p.Status = status[1]
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}

		var h []model.Holdings
		if err = config.DB.Table("holdings").Where("user_id=? AND stock_name=?", p.UserId, res.StockName).Find(&h).Error; err != nil {
			return err
		}
		price := 0

		for _, check := range h {
			if res.Quantity >= check.Quantity {
				res.Quantity = res.Quantity - check.Quantity
				orderHist := model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      check.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     res.AveragePrice,
					CommissionFee: 2000,
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err = config.DB.Table("order_history").Create(orderHist).Error; err != nil {
					return err
				}
				if err = config.DB.Table("holdings").Where("id", check.Id).Delete(check).Error; err != nil {
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			} else {
				orderHist := model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      res.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     res.AveragePrice,
					CommissionFee: 2000,
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err = config.DB.Table("order_history").Create(orderHist).Error; err != nil {
					return err
				}
				check.Quantity = check.Quantity - res.Quantity

				if err = config.DB.Table("holdings").Where("id", check.Id).Updates(&check).Error; err != nil {
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}

		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			return err
		}
		account.Balance = account.Balance + int64(price)
		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			return err
		}

		if err = config.DB.Table("pending_orders").Where("order_id=?", res.OrderID).Delete(p).Error; err != nil {
			return err
		}

	} else if res.Status == status[2] {
		//Sell Order Half completed
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).First(&p).Error; err != nil {
			return err
		}
		p.Status = status[2]
		p.Quantity = p.Quantity - res.Quantity
		if err = config.DB.Table("pending_orders").Where("order_id", res.OrderID).Updates(&p).Error; err != nil {
			return err
		}

		var h []model.Holdings
		if err = config.DB.Table("holdings").Where("user_id=? AND stock_name=?", p.UserId, res.StockName).Find(&h).Error; err != nil {
			return err
		}
		price := 0
		for _, check := range h {
			if res.Quantity >= check.Quantity {
				res.Quantity = res.Quantity - check.Quantity
				orderHist := model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      check.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     res.AveragePrice,
					CommissionFee: 2000,
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err = config.DB.Table("order_history").Create(orderHist).Error; err != nil {
					return err
				}
				if err = config.DB.Table("holdings").Where("id", check.Id).Delete(check).Error; err != nil {
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			} else {
				orderHist := model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      res.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     res.AveragePrice,
					CommissionFee: 2000,
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err = config.DB.Table("order_history").Create(orderHist).Error; err != nil {
					return err
				}
				check.Quantity = check.Quantity - res.Quantity

				if err = config.DB.Table("holdings").Where("id", check.Id).Updates(&check).Error; err != nil {
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}
		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			return err
		}
		account.Balance = account.Balance + int64(price)
		if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			return err
		}
	}
	return nil
}

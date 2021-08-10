package order

import (
	_ "Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/services/stock_exchange"
	"log"
)

func UpdateBuyOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	var account model.TradingAccount

	if res.Status == Failed {
		//Buy order Failed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Failed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating status in pending orders", err)
			return err
		}
		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			log.Println("Error in fetching trading account", err)
			return err
		}

		if p.BookType == Market {
			account.Balance = account.Balance + int64(p.Quantity*p.OrderPrice)
		} else if p.BookType == Limit {
			account.Balance = account.Balance + int64(p.Quantity*p.LimitPrice)
		}
		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			log.Println("Error in updating balance in trading account", err)
			return err
		}
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Delete(&p).Error; err != nil {
			log.Println("Error in deleting failed order in pending orders", err)
			return err
		}
		return nil
	} else if res.Status == Completed {
		//Buy order completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating status in pending orders", err)
			return err
		}
		h := model.Holdings{
			UserId:    p.UserId,
			OrderId:   p.OrderId,
			StockName: p.StockName,
			Quantity:  p.Quantity,
			BuyPrice:  res.AveragePrice,
			OrderedAt: res.OrderExecutionTime,
		}
		if p.BookType == Market {
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				log.Println("Error in fetching trading account", err)
				return err
			}
			account.Balance = account.Balance + int64((p.OrderPrice-res.AveragePrice)*res.Quantity)
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				log.Println("Error in updating balance in trading account", err)
				return err
			}
		} else if p.BookType == Limit {
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				log.Println("Error in fetching trading account", err)
				return err
			}
			account.Balance = account.Balance + int64((p.LimitPrice-res.AveragePrice)*res.Quantity)
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				log.Println("Error in updating balance in trading account", err)
				return err
			}
		}

		if err = database.GetDB().Table("holdings").Create(&h).Error; err != nil {
			log.Println("Error in creating holdings", err)
			return err
		}
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Delete(&p).Error; err != nil {
			log.Println("Error in deleting order in pending orders", err)
			return err
		}
	} else if res.Status == Partial {
		//Buy order Half Completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Partial
		p.Quantity = p.Quantity - res.Quantity
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating orders in pending orders", err)
			return err
		}

		h := model.Holdings{
			UserId:    p.UserId,
			OrderId:   p.OrderId,
			StockName: p.StockName,
			Quantity:  res.Quantity,
			BuyPrice:  res.AveragePrice,
			OrderedAt: res.OrderExecutionTime,
		}
		if p.BookType == Market {
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				log.Println("Error in fetching trading account", err)
				return err
			}
			account.Balance = account.Balance + int64((p.OrderPrice-res.AveragePrice)*res.Quantity)
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				log.Println("Error in updating balance in trading account", err)
				return err
			}
		} else if p.BookType == Limit {
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				log.Println("Error in fetching trading account", err)
				return err
			}
			account.Balance = account.Balance + int64((p.LimitPrice-res.AveragePrice)*res.Quantity)
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				log.Println("Error in updating balance in trading account", err)
				return err
			}
		}
		if err = database.GetDB().Table("holdings").Create(&h).Error; err != nil {
			log.Println("Error in creating holdings", err)
			return err
		}
	}
	return nil
}

func UpdateSellOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	var account model.TradingAccount

	if res.Status == Failed {
		//Sell order failed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Failed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating status in pending orders", err)
			return err
		}
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Delete(&p).Error; err != nil {
			log.Println("Error in deleting order in pending orders", err)
			return err
		}
		return nil
	} else if res.Status == Completed {
		//Sell order Completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating order in pending orders", err)
			return err
		}

		var h []model.Holdings
		if err = database.GetDB().Table("holdings").Where("user_id=? AND stock_name=?", p.UserId, res.StockName).Find(&h).Error; err != nil {
			log.Println("Error in fetching holdings", err)
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
				}
				if err = database.GetDB().Table("order_history").Create(&orderHist).Error; err != nil {
					log.Println("Error in creating order history", err)
					return err
				}
				if err = database.GetDB().Table("holdings").Where("id=?", check.Id).Delete(&check).Error; err != nil {
					log.Println("Error in deleting holdings", err)
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
				}
				if err = database.GetDB().Table("order_history").Create(&orderHist).Error; err != nil {
					log.Println("Error in creating order history", err)
					return err
				}
				check.Quantity = check.Quantity - res.Quantity
				if err = database.GetDB().Table("holdings").Where("id=?", check.Id).Updates(&check).Error; err != nil {
					log.Println("Error in updating holdings", err)
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}

		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			log.Println("Error in fetching trading account", err)
			return err
		}
		account.Balance = account.Balance + int64(price)
		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			log.Println("Error in updating balance in trading account", err)
			return err
		}

		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Delete(&p).Error; err != nil {
			log.Println("Error in deleting order in pending orders", err)
			return err
		}

	} else if res.Status == Partial {
		//Sell Order Half completed
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).First(&p).Error; err != nil {
			log.Println("Error in fetching pending orders", err)
			return err
		}
		p.Status = Partial
		p.Quantity = p.Quantity - res.Quantity
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", res.OrderID).Updates(&p).Error; err != nil {
			log.Println("Error in updating order in pending orders", err)
			return err
		}

		var h []model.Holdings
		if err = database.GetDB().Table("holdings").Where("user_id=? AND stock_name=?", p.UserId, res.StockName).Find(&h).Error; err != nil {
			log.Println("Error in fetching holdings", err)
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
				}
				if err = database.GetDB().Table("order_history").Create(&orderHist).Error; err != nil {
					log.Println("Error in creating order history", err)
					return err
				}
				if err = database.GetDB().Table("holdings").Where("id=?", check.Id).Delete(&check).Error; err != nil {
					log.Println("Error in deleting holdings", err)
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
				}
				if err = database.GetDB().Table("order_history").Create(&orderHist).Error; err != nil {
					log.Println("Error in creating order history", err)
					return err
				}
				check.Quantity = check.Quantity - res.Quantity
				if err = database.GetDB().Table("holdings").Where("id=?", check.Id).Updates(&check).Error; err != nil {
					log.Println("Error in updating holdings", err)
					return err
				}
				price = price + orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}
		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
			log.Println("Error in fetching trading account", err)
			return err
		}
		account.Balance = account.Balance + int64(price)
		if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
			log.Println("Error in updating balance in trading account", err)
			return err
		}
	}
	return nil
}

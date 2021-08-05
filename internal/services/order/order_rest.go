package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"time"
)

func UpdateBuyOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	if res.Status=="Failed"{
		p.Status="Failed"
		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Delete(p).Error; err!=nil{
			return err
		}
		return nil

	}else if res.Status=="Completed" {
		p.Status="Completed"
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}
		h:=model.Holdings{
			UserId: p.UserId,
			OrderId: p.OrderId,
			StockName: p.StockName,
			Quantity: p.Quantity,
			BuyPrice: int(res.AveragePrice),
			OrderedAt: p.CreatedAt,
			CreatedAt: time.Now(),
			UpdatedAt: res.OrderExecutionTime,
		}
		if err=config.DB.Table("holdings").Create(h).Error;err!=nil{
			return err
		}
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Delete(p).Error; err!=nil{
			return err
		}
	} else if res.Status=="Half Completed" {
		p.Status="Half Completed"
		p.Quantity=p.Quantity-int(res.Quantity)
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}

		h:=model.Holdings{
			UserId: p.UserId,
			OrderId: p.OrderId,
			StockName: p.StockName,
			Quantity: int(res.Quantity),
			BuyPrice: int(res.AveragePrice),
			OrderedAt: p.CreatedAt,
			CreatedAt: time.Now(),
			UpdatedAt: res.OrderExecutionTime,
		}
		if err=config.DB.Table("holdings").Create(h).Error;err!=nil{
			return err
		}
	}
	return nil
}

func UpdateSellOrder(res *stock_exchange.OrderResponse) (err error){
	var p model.PendingOrders
	if res.Status=="failed"{
		p.Status="Failed"
		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Delete(p).Error; err!=nil{
			return err
		}
		return nil
	} else if res.Status=="Completed"{
		p.Status="Completed"
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}

		var h []model.Holdings
		if err=config.DB.Table("holdings").Where("user_id=? AND stock_name=?",p.UserId,res.StockName).Find(&h).Error;err!=nil{
			return err
		}
		for _,check:=range h{
			if int(res.Quantity)>=check.Quantity{
				res.Quantity=res.Quantity-uint(check.Quantity)
				orderHist:=model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      check.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     int(res.AveragePrice),
					CommissionFee: int(float64(check.Quantity*int(res.AveragePrice))*0.01),
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err=config.DB.Table("order_history").Create(orderHist).Error;err!=nil{
					return err
				}
				if err=config.DB.Table("holdings").Where("id",check.Id).Delete(check).Error;err!=nil{
					return err
				}
			} else {
				orderHist:=model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      int(res.Quantity),
					BuyPrice:      check.BuyPrice,
					SellPrice:     int(res.AveragePrice),
					CommissionFee: int(float64(res.Quantity*res.AveragePrice)*0.01),
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err=config.DB.Table("order_history").Create(orderHist).Error;err!=nil{
					return err
				}
				check.Quantity=check.Quantity-int(res.Quantity)

				if err=config.DB.Table("holdings").Where("id",check.Id).Update(check).Error;err!=nil{
					return err
				}
			}
		}
		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Delete(p).Error;err!=nil{
			return err
		}
	} else if res.Status=="Half Completed"{
		p.Status="Half Completed"
		p.Quantity=p.Quantity-int(res.Quantity)
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}

		var h []model.Holdings
		if err=config.DB.Table("holdings").Where("user_id=? AND stock_name=?",p.UserId,res.StockName).Find(&h).Error;err!=nil{
			return err
		}
		for _,check:=range h{
			if int(res.Quantity)>=check.Quantity{
				res.Quantity=res.Quantity-uint(check.Quantity)
				orderHist:=model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      check.Quantity,
					BuyPrice:      check.BuyPrice,
					SellPrice:     int(res.AveragePrice),
					CommissionFee: int(float64(check.Quantity*int(res.AveragePrice))*0.01),
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err=config.DB.Table("order_history").Create(orderHist).Error;err!=nil{
					return err
				}
				if err=config.DB.Table("holdings").Where("id",check.Id).Delete(check).Error;err!=nil{
					return err
				}
			} else {
				orderHist:=model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      int(res.Quantity),
					BuyPrice:      check.BuyPrice,
					SellPrice:     int(res.AveragePrice),
					CommissionFee: int(float64(res.Quantity*res.AveragePrice)*0.01),
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err=config.DB.Table("order_history").Create(orderHist).Error;err!=nil{
					return err
				}
				check.Quantity=check.Quantity-int(res.Quantity)

				if err=config.DB.Table("holdings").Where("id",check.Id).Update(check).Error;err!=nil{
					return err
				}
			}
		}
	}
	return nil
}

package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"time"
)

func UpdateBuyOrder(res *stock_exchange.OrderResponse) (err error) {
	var p model.PendingOrders
	var balance model.Payments
	if res.Status=="FAILED"{
		p.Status="FAILED"
		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Updates(&p).Error; err!=nil{
			return err
		}
		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).First(&balance).Error;err!=nil{
			return err
		}
		balance.CurrentBalance=balance.CurrentBalance+p.Quantity*p.OrderPrice
		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).Updates(&balance).Error;err!=nil{
			return err
		}
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Delete(p).Error; err!=nil{
			return err
		}
		return nil

	}else if res.Status=="COMPLETED" {
		p.Status="COMPLETED"
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(&p).Error; err!=nil{
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
	} else if res.Status=="HALF_COMPLETED" {
		p.Status="HALF_COMPLETED"
		p.Quantity=p.Quantity-int(res.Quantity)
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(&p).Error; err!=nil{
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
	var balance model.Payments

	if res.Status=="FAILED"{
		p.Status="FAILED"
		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Updates(&p).Error; err!=nil{
			return err
		}
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Delete(p).Error; err!=nil{
			return err
		}
		return nil
	} else if res.Status=="COMPLETED"{
		p.Status="COMPLETED"
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(&p).Error; err!=nil{
			return err
		}

		var h []model.Holdings
		if err=config.DB.Table("holdings").Where("user_id=? AND stock_name=?",p.UserId,res.StockName).Find(&h).Error;err!=nil{
			return err
		}
		price:=0

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
					CommissionFee: 2000,
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
				price= price+ orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			} else {
				orderHist:=model.OrderHistory{
					UserId:        check.UserId,
					OrderId:       check.OrderId,
					StockName:     check.StockName,
					Quantity:      int(res.Quantity),
					BuyPrice:      check.BuyPrice,
					SellPrice:     int(res.AveragePrice),
					CommissionFee: 2000,
					BoughtAt:      check.OrderedAt,
					SoldAt:        res.OrderExecutionTime,
					CreatedAt:     time.Now(),
				}
				if err=config.DB.Table("order_history").Create(orderHist).Error;err!=nil{
					return err
				}
				check.Quantity=check.Quantity-int(res.Quantity)

				if err=config.DB.Table("holdings").Where("id",check.Id).Updates(check).Error;err!=nil{
					return err
				}
				price= price+ orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}

		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).First(&balance).Error;err!=nil{
			return err
		}
		balance.CurrentBalance=balance.CurrentBalance + price
		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).Updates(&balance).Error;err!=nil{
			return err
		}

		if err=config.DB.Table("pending_orders").Where("order_id=?",res.OrderID).Delete(p).Error;err!=nil{
			return err
		}

	} else if res.Status=="HALF_COMPLETED"{
		p.Status="HALF_COMPLETED"
		p.Quantity=p.Quantity-int(res.Quantity)
		if err=config.DB.Table("pending_orders").Where("order_id",res.OrderID).Updates(p).Error; err!=nil{
			return err
		}

		var h []model.Holdings
		if err=config.DB.Table("holdings").Where("user_id=? AND stock_name=?",p.UserId,res.StockName).Find(&h).Error;err!=nil{
			return err
		}
		price:=0
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
				price= price+ orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
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

				if err=config.DB.Table("holdings").Where("id",check.Id).Updates(&check).Error;err!=nil{
					return err
				}
				price= price+ orderHist.Quantity*orderHist.SellPrice - orderHist.CommissionFee
			}
		}
		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).First(&balance).Error;err!=nil{
			return err
		}
		balance.CurrentBalance=balance.CurrentBalance + price
		if err=config.DB.Table("payments").Where("user_id=?",p.UserId).Updates(&balance).Error;err!=nil{
			return err
		}


	}
	return nil
}

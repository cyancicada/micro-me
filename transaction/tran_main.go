package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type (
	User struct {
		Id       int64
		Username string `json:"username" xorm:"varchar(100) notnull 'username'"`
		Password string `json:"password" xorm:"varchar(100) notnull 'password'"`
	}
	Balance struct {
		Id      int64
		UserId  int64   `json:"userId" xorm:"int notnull 'user_id'"`
		Balance float64 `json:"balance" xorm:"double notnull 'balance'"`
	}

	Goods struct {
		Id        int64
		GoodsName string  `json:"goodsName" xorm:"varchar(255) notnull 'goods_name'"`
		Price     float64 `json:"price" xorm:"double notnull 'price'"`
		Stock     int64   `json:"stock" xorm:"int notnull 'stock'"`
	}
	Order struct {
		Id     int64
		UserId int64   `json:"userId" xorm:"int notnull 'user_id'"`
		Amount float64 `json:"amount" xorm:"double notnull 'amount'"`
	}
)

var engineUser *xorm.Engine
var engineOrder *xorm.Engine
var enginebalance *xorm.Engine
var engineGoods *xorm.Engine
var err error

func main() {
	engineUser, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.5.100:3306)/members?charset=utf8")
	engineOrder, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.5.100:3306)/members?charset=utf8")
	enginebalance, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.5.100:3307)/balance?charset=utf8")
	engineGoods, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.5.100:3308)/goods?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	user := new(User)
	order := new(Order)
	balance := new(Balance)
	goods := new(Goods)
	if _, err := engineUser.Where("id = ?", 1).Get(user); err != nil {
		log.Fatal(err)
	}

	//订单服务
	orderNum := float64(2)
	goodsId := int64(1)
	one := goods.findById(goodsId)
	amount := orderNum * one.Price
	fmt.Println(one)
	err = order.insertOrderRecord(user.Id, amount, func() error {

		//订单业务代码1，
		//订单业务代码2;
		//订单业务代码2n;
		return balance.updateBalanceByUserId(user.Id, amount, func() error {
			//  余额业务代码1，
			//  余额业务代码2;
			//   余额业务代码2n;
			return goods.updateStockByGoodsId(one.Id, one.Stock-2, func() error {
				//商品业务代码1
				//商品业务代码2
				//商品业务代码n
				return nil
			})
		})
		//订单业务代码2n;
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success")

}

func (b *Balance) updateBalanceByUserId(userId int64, amount float64, opt func() error) error {
	_, err := enginebalance.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		one := b.findBalanceByUserId(userId)
		fmt.Println(amount)
		if _, err := session.Where("user_id = ?", userId).Update(&Balance{Balance: one.Balance - amount}); err != nil {
			log.Println("updateBalanceByUserId", err)
			return nil, err
		}
		return nil, opt()
	})
	return err
}
func (b *Balance) findBalanceByUserId(userId int64) *Balance {
	one := new(Balance)
	if _, err := enginebalance.Where("user_id = ?", userId).Get(one); err != nil {
		return nil
	}
	return one

}
func (g *Goods) updateStockByGoodsId(godsId int64, stock int64, opt func() error) error {
	_, err := engineGoods.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		fmt.Println("stock", stock)
		if _, err := session.Where("id = ?", godsId).Update(&Goods{Stock: stock}); err != nil {
			log.Println("updateStockByGoodsId", err)
			return nil, err
		}
		return nil, opt()
	})
	return err
}

func (g *Goods) findById(godsId int64) *Goods {
	goods := new(Goods)
	if _, err := engineGoods.Where("id = ?", godsId).Get(goods); err != nil {
		return nil
	}
	return goods

}
func (o *Order) insertOrderRecord(userId int64, amount float64, opt func() error) error {
	_, err := engineOrder.Transaction(func(session *xorm.Session) (i interface{}, e error) {

		if _, err := session.Insert(&Order{Amount: amount, UserId: userId}); err != nil {
			log.Println("insertOrderRecord", err)
			return nil, err
		}
		return nil, opt()
	})
	return err
}

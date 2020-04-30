package test

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
)

var db *dbx.Database
func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		Host:       "192.168.43.171:3306",
		User:       "root",
		Password:   "root",
		Database:   "test",
		Options: map[string]string{
			"parseTime": "true",
		},
	}
	var err error
	db, err = dbx.Open(settings)
	if err != nil {
		fmt.Println(err)
	}
	db.RegisterTable(&GoodsSigned{}, "goods")
	db.RegisterTable(&GoodsUnsigned{}, "goods_unsigned")
}

var query = "select * from goods" +
	" where envelope_no=?" +
	" for update"
var update = "update goods " +
	"set remain_amount=?,remain_quantity=?" +
	" where envelope_no=?"

func UpdateForLock(g *Goods) {
	err := db.Tx(func (runner *dbx.TxRunner) error {
		out := &GoodsSigned{}
		ok, err := runner.Get(out, query, g.EnvelopeNo)
		if !ok || err != nil {
			return err
		}

		subAmount := decimal.NewFromFloat(0.01)
		remainAmount := out.RemainAmount.Sub(subAmount)
		remainQuantity := out.RemainQuantity - 1

		_, row, err := runner.Execute(update, remainAmount, remainQuantity, g.EnvelopeNo)
		if err != nil {
			return err
		}

		if row < 1 {
			return errors.New("库存扣减失败")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func UpdateForUnsigned(g *Goods) {
	update := "update goods_unsigned " +
		"set remain_amount=remain_amount-?," +
		"remain_quantity=remain_quantity-? " +
		"where envelope_no=?"
	_, row, err := db.Execute(update, 0.01, 1, g.EnvelopeNo)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}

//乐观锁方案
func UpdateForOptimistic(g *Goods) {
	update := "update goods " +
		"set remain_amount=remain_amount-?, " +
		" remain_quantity=remain_quantity-? " +
		" where envelope_no=? " +
		" and remain_amount>=? " +
		" and remain_quantity>=? "
	_, row, err := db.Execute(update, 0.01, 1, g.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}


//乐观锁+无符号字段双保险方案
func UpdateForOptimisticAndUnsigned(g *Goods) {
	update := "update goods_unsigned " +
		"set remain_amount=remain_amount-?, " +
		" remain_quantity=remain_quantity-? " +
		" where envelope_no=? " +
		" and remain_amount>=? " +
		" and remain_quantity>=? "
	_, row, err := db.Execute(update, 0.01, 1, g.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if row < 1 {
		fmt.Println("扣减失败")
	}
}

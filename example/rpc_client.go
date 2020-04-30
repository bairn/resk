package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"resk/services"
)

func main() {
	c, err := rpc.Dial("tcp", ":18082")
	if err != nil {
		logrus.Panic(err)
	}
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "1bARrIz6IEVGvrJgIMArNC631ky",
		Username:     "测试用户",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	err = c.Call("EnvelopeRpc.SendOut", in, &out)
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("%+v", out)
}
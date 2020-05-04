package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bairn/account/core/accounts"
	"github.com/bairn/account/services"
	"github.com/bairn/infra/algo"
	"github.com/bairn/infra/base"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	envelopeServices "resk/services"
)

var multiple = decimal.NewFromFloat(100.0)

func (d *goodsDomain) Receive(ctx context.Context, dto envelopeServices.RedEnvelopeReceiveDTO) (item *envelopeServices.RedEnvelopeItemDTO, err error) {
	d.preCreateItem(dto)
	goods := d.Get(dto.EnvelopeNo)
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <=0 {
		log.Errorf("%+v", goods)
		return nil, errors.New("没有足够的红包和金额了")
	}
	nextAmount := d.nextAmount(goods)
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner:runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)

		if rows <= 0 || err != nil {
			log.Errorf("rows=%d, %s", rows, err.Error())
			return errors.New("没有足够的红包和金额了")
		}

		d.item.Quantity = 1
		d.item.PayStatus = int(envelopeServices.Paying)
		d.item.AccountNo = dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount

		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.item.Save(txCtx)
		if err != nil {
			log.Error(err)
			return err
		}

		status, err := d.transfer(txCtx, dto)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		return err
	})
	return d.item.ToDTO(), err
}

func (d *goodsDomain) transfer(ctx context.Context, dto envelopeServices.RedEnvelopeReceiveDTO) (status services.TransferedStatus, err error) {
	systemAmount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAmount.AccountNo,
		UserId:    systemAmount.UserId,
		Username:  systemAmount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}

	transfer := services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Desc:        "红包收入",
	}
	adomain := accounts.NewAccountDomain()
	return adomain.TransferWithContextTx(ctx, transfer)
}

func (d *goodsDomain) preCreateItem(dto envelopeServices.RedEnvelopeReceiveDTO) {
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{String: dto.RecvUsername, Valid: true}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}

	if goods.EnvelopeType == envelopeServices.GeneralEnvelopeType {
		return goods.AmountOne
	} else if goods.EnvelopeType == envelopeServices.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	} else {
		log.Error("不支持的红包类型")
	}
	return amount
}
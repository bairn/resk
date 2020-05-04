package envelopes

import (
	"context"
	"errors"
	"github.com/bairn/account/services"
	"github.com/bairn/infra/base"
	envelopeServices "github.com/bairn/resk/services"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		envelopeServices.IRedEnvelopeService = new(redEnvelopeService)
	})
}

type redEnvelopeService struct {
}

//发红包
func (r *redEnvelopeService) SendOut(dto envelopeServices.RedEnvelopeSendingDTO) (activity *envelopeServices.RedEnvelopeActivity, err error) {
	//验证
	if err = base.ValidateStruct(&dto); err != nil {
		return activity, err
	}
	//获取红包发送人的资金账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在：" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo

	if goods.Blessing == "" {
		goods.Blessing = envelopeServices.DefaultBlessing
	}
	if goods.EnvelopeType == envelopeServices.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	//执行发送红包的逻辑
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		log.Error(err)
	}

	return activity, err
}

func (r *redEnvelopeService) Receive(dto envelopeServices.RedEnvelopeReceiveDTO) (item *envelopeServices.RedEnvelopeItemDTO, err error) {
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}

	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包资金账户不存在:user_id=" + dto.RecvUserId)
	}
	dto.AccountNo = account.AccountNo

	domain := goodsDomain{}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}

func (r *redEnvelopeService) Refund(envelopeNo string) (order *envelopeServices.RedEnvelopeGoodsDTO) {
	panic("implement me")
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *envelopeServices.RedEnvelopeGoodsDTO) {
	domain := goodsDomain{}
	po := domain.GetOne(envelopeNo)
	if po == nil {
		return order
	}
	return po.ToDTO()
}


func (r *redEnvelopeService) ListSent(userId string, page, size int) (orders []*envelopeServices.RedEnvelopeGoodsDTO) {
	domain := new(goodsDomain)
	pos := domain.FindByUser(userId, page, size)
	orders = make([]*envelopeServices.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, p := range pos {
		orders = append(orders, p.ToDTO())
	}

	return
}

func (r *redEnvelopeService) ListReceivable(page, size int) (orders []*envelopeServices.RedEnvelopeGoodsDTO) {
	domain := new(goodsDomain)

	pos := domain.ListReceivable(page, size)
	orders = make([]*envelopeServices.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, p := range pos {
		if p.RemainQuantity > 0 {
			orders = append(orders, p.ToDTO())
		}
	}
	return
}

func (r *redEnvelopeService) ListReceived(userId string, page, size int) (items []*envelopeServices.RedEnvelopeItemDTO) {
	domain := new(goodsDomain)
	pos := domain.ListReceived(userId, page, size)
	items = make([]*envelopeServices.RedEnvelopeItemDTO, 0, len(pos))
	for _, p := range pos {
		items = append(items, p.ToDTO())
	}
	return
}

func (r *redEnvelopeService) ListItems(envelopeNo string) (items []*envelopeServices.RedEnvelopeItemDTO) {
	domain := itemDomain{}
	return domain.FindItems(envelopeNo)
}
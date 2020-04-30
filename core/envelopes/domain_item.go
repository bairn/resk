package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
)

type itemDomain struct {
	RedEnvelopeItem
}

func (d *itemDomain) createItemNo() {
	d.ItemNo = ksuid.New().Next().String()
}

func (d *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	d.RedEnvelopeItem.FromDTO(&item)
	d.RecvUsername.Valid = true
	d.createItemNo()
}

func (d *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner:runner}
		id, err = dao.Insert(&d.RedEnvelopeItem)
		return err
	})
	return id, err
}

func (d *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner:runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

func (d *itemDomain) FindItems(envelopeNo string) (itemDtos []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner:runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return itemDtos
	}

	itemDtos = make([]*services.RedEnvelopeItemDTO, 0)
	for _, po := range items {
		itemDtos = append(itemDtos, po.ToDTO())
	}
	return itemDtos
}

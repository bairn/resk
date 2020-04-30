package test

import (
	"github.com/shopspring/decimal"
	"time"
)

type GoodsSigned struct {
	Goods
}

type GoodsUnsigned struct {
	Goods
}

type Goods struct {
	RemainAmount decimal.Decimal `db:"remain_amount"`
	RemainQuantity int `db:"remain_quantity"`
	EnvelopeNo string `db:"envelope_no,uni"`
	CreateAt time.Time `db:"created_at,omitempty"`
	UpdateAt time.Time `db:"updated_at,omitempty"`
}

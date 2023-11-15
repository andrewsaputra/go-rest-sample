package internal

import "github.com/rs/xid"

func NewXidGenerator() *XidGenerator {
	return &XidGenerator{}
}

type XidGenerator struct{}

func (this *XidGenerator) NextId() string {
	return xid.New().String()
}

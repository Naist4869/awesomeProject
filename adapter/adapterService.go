package adapter

import (
	"context"

	"github.com/Naist4869/awesomeProject/model/jdmodel"
)

type Service interface {
	CategoryGoodsGet(context.Context, jdmodel.CategoryGoodsGetReq) (jdmodel.CategoryGoodsGet, error)
}

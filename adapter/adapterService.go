package adapter

import "context"

type Adapter interface {
	CategoryGoodsGet(ctx context.Context)
}

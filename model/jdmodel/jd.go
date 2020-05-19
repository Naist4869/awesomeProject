package jdmodel

type CategoryGoodsGetReq struct {
	ParentId int64 `json:"parentId"` // 父类目id(一级父类目为0)
	Grade    int64 `json:"grade"`    // 类目级别(类目级别 0，1，2 代表一、二、三级类目)
}

type CategoryGoodsGet struct {
	ID       int    `json:"id"`       // 类目Id
	Name     string `json:"name"`     // 类目名称
	Grade    int    `json:"grade"`    // 类目级别(类目级别 0，1，2 代表一、二、三级类目)
	ParentId int    `json:"parentId"` // 父类目Id
}

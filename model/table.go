package model

import "errors"

type Table struct {
	Query map[string]interface{} `json:"query"` // 查询条件
	Sort  []string               `json:"sort"`  // 排序
	Start int                    `json:"start"` // 开始位置，从0开始
	Limit int                    `json:"limit"` // 本次最多数量
}

func (t Table) Validate() error {

	if t.Limit < 0 {
		return errors.New("limit必须大于等于0")
	}
	if t.Start < 0 {
		return errors.New("start必须大于等于0")
	}
	return nil

}

package usermodel

// 用户状态
const (
	_         = iota
	NotActive // 未激活
	Normal    // 正常(激活)
	Freeze    // 冻结
)

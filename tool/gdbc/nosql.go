package gdbc

// DBComponent 对需要数据库的业务模块的抽象
type DBComponent interface {
	Keys() map[string]*Spec
	Init() error
}

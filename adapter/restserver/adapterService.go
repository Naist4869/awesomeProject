package restserver

type Router interface {
	Router() IRouter
}

// IRouter http router framework interface.
type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

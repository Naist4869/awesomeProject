package restserver

func monitor() HandlerFunc {
	return func(c *Context) {
		h := metrics.Handler()
		h.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

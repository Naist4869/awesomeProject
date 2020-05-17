package restserver

type options struct {
	// If enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool

	// If true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool
}

type ServerOption interface {
	applyTo(*options)
}

func defaultOptions() options {
	return options{}
}

type OptionUseRawPath struct {
	x bool
}

var _ ServerOption = (*OptionUseRawPath)(nil)

func (x OptionUseRawPath) applyTo(y *options) {
	y.UseRawPath = x.x
}

type OptionUnescapePathValues struct {
	x bool
}

var _ ServerOption = (*OptionUseRawPath)(nil)

func (x OptionUnescapePathValues) applyTo(y *options) {
	y.UnescapePathValues = x.x
}

type OptionHandleMethodNotAllowed struct {
	x bool
}

var _ ServerOption = (*OptionHandleMethodNotAllowed)(nil)

func (x OptionHandleMethodNotAllowed) applyTo(y *options) {
	y.HandleMethodNotAllowed = x.x
}

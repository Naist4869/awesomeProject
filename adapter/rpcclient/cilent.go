package rpcclient

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"github.com/Naist4869/log"
	"go.uber.org/zap"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

var (
	_once           sync.Once
	_defaultCliConf = &ClientConfig{
		Dial:              time.Second * 10,
		Timeout:           time.Second * 5,
		Subset:            50,
		KeepAliveInterval: time.Second * 60,
		KeepAliveTimeout:  time.Second * 20,
	}
	_defaultClient *Client
	_abortIndex    int8 = math.MaxInt8 / 2
)

// ClientConfig is rpc client conf.
type ClientConfig struct {
	Dial                   time.Duration
	Timeout                time.Duration
	Method                 map[string]*ClientConfig
	Clusters               []string
	Zone                   string
	Subset                 int
	NonBlock               bool
	KeepAliveInterval      time.Duration
	KeepAliveTimeout       time.Duration
	KeepAliveWithoutStream bool
}

// SetConfig hot reloads client config
func (c *Client) SetConfig(conf *ClientConfig) (err error) {
	if conf == nil {
		conf = _defaultCliConf
	}
	if conf.Dial <= 0 {
		conf.Dial = time.Second * 10
	}
	if conf.Timeout <= 0 {
		conf.Timeout = time.Millisecond * 250
	}
	if conf.Subset <= 0 {
		conf.Subset = 50
	}
	if conf.KeepAliveInterval <= 0 {
		conf.KeepAliveInterval = time.Second * 60
	}
	if conf.KeepAliveTimeout <= 0 {
		conf.KeepAliveTimeout = time.Second * 20
	}

	c.mutex.Lock()
	c.conf = conf
	c.mutex.Unlock()
	return
}

// TimeoutCallOption timeout option.
type TimeoutCallOption struct {
	*grpc.EmptyCallOption
	Timeout time.Duration
}

// WithTimeoutCallOption can override the timeout in ctx and the timeout in the configuration file
func WithTimeoutCallOption(timeout time.Duration) *TimeoutCallOption {
	if timeout <= 0 {
		timeout = _defaultCliConf.Timeout
	}
	return &TimeoutCallOption{&grpc.EmptyCallOption{}, timeout}
}

// Client is the framework's client side instance, it contains the ctx, opt and interceptors.
// Create an instance of Client, by using NewClient().
type Client struct {
	conf   *ClientConfig
	mutex  sync.RWMutex
	logger *log.Logger

	opts     []grpc.DialOption
	handlers []grpc.UnaryClientInterceptor
}

// UseOpt attachs a global grpc DialOption to the Client.
func (c *Client) UseOpt(opts ...grpc.DialOption) *Client {
	c.opts = append(c.opts, opts...)
	return c
}
func (c *Client) cloneOpts() []grpc.DialOption {
	dialOptions := make([]grpc.DialOption, len(c.opts))
	copy(dialOptions, c.opts)
	return dialOptions
}

// NewClient returns a new blank Client instance with a default client interceptor.
// opt can be used to add grpc dial options.
func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *Client {
	c := new(Client)
	if err := c.SetConfig(conf); err != nil {
		panic(err)
	}
	//c.UseOpt(grpc.WithBalancerName(p2c.Name))
	c.UseOpt(opt...)
	return c
}

// DefaultClient returns a new default Client instance with a default client interceptor and default dialoption.
// opt can be used to add grpc dial options.
func DefaultClient() *Client {
	_once.Do(func() {
		_defaultClient = NewClient(nil)
		_defaultClient.logger = log.BaseLogger
	})
	return _defaultClient
}

func (c *Client) handle() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		var (
			ok      bool
			conf    *ClientConfig
			cancel  context.CancelFunc
			timeOpt *TimeoutCallOption
		)
		c.mutex.RLock()
		if conf, ok = c.conf.Method[method]; !ok {
			conf = c.conf
		}
		c.mutex.RUnlock()
		for _, opt := range opts {
			var tok bool
			timeOpt, tok = opt.(*TimeoutCallOption)
			if tok {
				break
			}
		}
		if timeOpt != nil && timeOpt.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), timeOpt.Timeout)
		} else {
			d := conf.Timeout
			// 比较全局配置和方法配置timeout换成两者之间较少的timeout
			if deadline, ok := ctx.Deadline(); ok {
				if ctimeout := time.Until(deadline); ctimeout > d {
					ctx, cancel = context.WithTimeout(ctx, d)
				}
			}
		}
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()
		if err = invoker(ctx, method, req, reply, cc, opts...); err != nil {
			return
		}

		return
	}
}

// NewConn will create a grpc conn by default config.
func NewConn(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultClient().Dial(context.Background(), target, opt...)
}

// Dial creates a client connection to the given target.
// Target format is scheme://authority/endpoint?query_arg=value
// example: discovery://default/account.account.service?cluster=shfy01&cluster=shfy02
func (c *Client) Dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts, grpc.WithInsecure())
	return c.dial(ctx, target, opts...)
}

// DialTLS creates a client connection over tls transport to the given target.
func (c *Client) DialTLS(ctx context.Context, target string, file string, name string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	return c.dial(ctx, target, opts...)
}
func (c *Client) dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	dialOptions := c.cloneOpts()
	if !c.conf.NonBlock {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                c.conf.KeepAliveInterval,
		Timeout:             c.conf.KeepAliveTimeout,
		PermitWithoutStream: !c.conf.KeepAliveWithoutStream,
	}))
	dialOptions = append(dialOptions, opts...)
	var handlers []grpc.UnaryClientInterceptor
	handlers = append(handlers, c.recovery())
	handlers = append(handlers, Logger(log.BaseLogger))
	handlers = append(handlers, c.handlers...)
	// NOTE: c.handle must be a last interceptor.
	handlers = append(handlers, c.handle())
	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(chainUnaryClient(handlers)))
	c.mutex.RLock()
	conf := c.conf
	c.mutex.RUnlock()
	if conf.Dial > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, conf.Dial)
		defer cancel()
	}
	// todo 集群

	if conn, err = grpc.DialContext(ctx, target, dialOptions...); err != nil {
		c.logger.Error("dial", zap.Error(err), zap.Stack("栈"))
	}
	return
}

// Use attachs a global inteceptor to the Client.
// For example, this is the right place for a circuit breaker or error management inteceptor.
func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	finalSize := len(c.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: client use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryClientInterceptor, finalSize)
	copy(mergedHandlers, c.handlers)
	copy(mergedHandlers[len(c.handlers):], handlers)
	c.handlers = mergedHandlers
	return c
}

// chainUnaryClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func chainUnaryClient(handlers []grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(handlers)
	if n == 0 {
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)
		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, ic *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == n-1 {
				return invoker(ictx, imethod, ireq, ireply, ic, iopts...)
			}
			i++
			return handlers[i](ictx, imethod, ireq, ireply, ic, chainHandler, iopts...)
		}

		return handlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}

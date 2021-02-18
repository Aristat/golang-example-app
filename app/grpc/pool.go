package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

var (
	poolList map[string]*Pool
	mu       sync.Mutex
)

// Done
type Done func()

type conn struct {
	conn *grpc.ClientConn
	pool *Pool
	err  error
}

func (c *conn) init() {
	c.conn, c.err = grpc.Dial(c.pool.target, c.pool.opts.dialOptions...)
}

// Pool
type Pool struct {
	id      string
	ctx     context.Context
	service string
	pool    *grpcpool.Pool
	target  string
	opts    *opts
}

// Get
func (p *Pool) Get() (*grpc.ClientConn, Done, error) {
	c, e := p.pool.Get(p.ctx)
	return c.ClientConn, func() {
		_ = c.Close()
	}, errors.WithMessage(e, prefix)
}

type opts struct {
	dialOptions     []grpc.DialOption
	initConn        int
	maxConn         int
	idleTimeout     time.Duration
	maxLifeDuration time.Duration
}

// Option
type Option func(*opts) error

// ConnOptions
func ConnOptions(o ...grpc.DialOption) Option {
	return func(f *opts) error {
		f.dialOptions = o
		return nil
	}
}

// MaxConn
func MaxConn(value int) Option {
	return func(f *opts) error {
		f.maxConn = value
		return nil
	}
}

// InitConn
func InitConn(value int) Option {
	return func(f *opts) error {
		f.initConn = value
		return nil
	}
}

// MaxLifeDuration
func MaxLifeDuration(value time.Duration) Option {
	return func(f *opts) error {
		f.maxLifeDuration = value
		return nil
	}
}

// IdleTimeout
func IdleTimeout(value time.Duration) Option {
	return func(f *opts) error {
		f.idleTimeout = value
		return nil
	}
}

// NewPool
func NewPool(ctx context.Context, service, target string, o ...Option) (_ *Pool, loaded bool) {
	mu.Lock()
	defer mu.Unlock()
	if p := poolList[service]; p != nil {
		return p, true
	}
	p := &Pool{
		id:      time.Now().String(),
		ctx:     ctx,
		service: service,
		opts:    &opts{},
		target:  target,
	}
	for _, option := range o {
		_ = option(p.opts)
	}
	factory := func() (*grpc.ClientConn, error) {
		c := &conn{pool: p}
		c.init()
		return c.conn, c.err
	}
	p.pool, _ = grpcpool.New(factory, p.opts.initConn, p.opts.maxConn, p.opts.idleTimeout, p.opts.maxLifeDuration)
	return p, false
}

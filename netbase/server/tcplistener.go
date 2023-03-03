package server

import (
	"crypto/tls"
	"github.com/lngwu11/toolgo/utils"
	"github.com/pkg/errors"
	"net"
	"time"
)

const (
	// 默认保活探测间隔
	defaultKeepAlivePeriod = 30 * time.Second
	// 默认network
	defaultNetwork = "tcp"
)

type Option func(opts *Options)

type Options struct {
	// 拨号配置
	network string
	// 保活探测间隔
	keepAlivePeriod time.Duration
	// tls配置
	tlsConfig *tls.Config
}

// WithNetwork .
func WithNetwork(network string) Option {
	return func(opts *Options) {
		opts.network = network
	}
}

// WithKeepAlivePeriod 。
func WithKeepAlivePeriod(keepAlivePeriod time.Duration) Option {
	return func(opts *Options) {
		opts.keepAlivePeriod = keepAlivePeriod
	}
}

// WithTLSConfig 。
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(opts *Options) {
		opts.tlsConfig = tlsConfig
	}
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	if len(opts.network) == 0 {
		opts.network = defaultNetwork
	}
	if opts.keepAlivePeriod == 0 {
		opts.keepAlivePeriod = defaultKeepAlivePeriod
	}
	return opts
}

func RunTCPListener(endpoint string, connCh chan<- utils.ConnReadWriteCloser, options ...Option) (listener net.Listener, err error) {
	opts := loadOptions(options...)

	switch opts.network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, net.UnknownNetworkError(opts.network)
	}

	if opts.tlsConfig != nil {
		// tls listen
		listener, err = tls.Listen(opts.network, endpoint, opts.tlsConfig)
		if err != nil {
			err = errors.Wrapf(err, "无法绑定tcp地址 [address=%v]", endpoint)
			return
		}
	} else {
		// normal listen
		listener, err = net.Listen(opts.network, endpoint)
		if err != nil {
			err = errors.Wrapf(err, "无法绑定tcp地址 [address=%v]", endpoint)
			return
		}
	}

	go func() {
		defer func() {
			_ = listener.Close()
		}()

		for {
			var c net.Conn
			c, err = listener.Accept()
			if err != nil {
				err = errors.Wrapf(err, "无法获取tcp连接")
				return
			}

			if tcp, ok := c.(*net.TCPConn); ok {
				_ = tcp.SetKeepAlive(true)
				_ = tcp.SetKeepAlivePeriod(opts.keepAlivePeriod)
			}

			connCh <- c.(utils.ConnReadWriteCloser)
		}
	}()

	return
}

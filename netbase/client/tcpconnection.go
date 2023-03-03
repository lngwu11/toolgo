package client

import (
	"crypto/tls"
	"github.com/lngwu11/toolgo/utils"
	"github.com/pkg/errors"
	"net"
	"time"
)

const (
	// 默认拨号超时时间
	defaultDialTimeout = 15 * time.Second
	// 默认保活探测时间
	defaultDialKeepAlive = 15 * time.Second
	// 默认连接空闲超时时间
	defaultConnIdleTimout = 60 * time.Second
	// 默认network
	defaultNetwork = "tcp"
)

type Option func(opts *Options)

type Options struct {
	// 空闲超时时间
	idleTimeout time.Duration
	// 拨号配置
	network       string
	dialTimeout   time.Duration
	dialKeepAlive time.Duration
	// tls配置
	tlsConfig *tls.Config
}

// WithNetwork .
func WithNetwork(network string) Option {
	return func(opts *Options) {
		opts.network = network
	}
}

// WithTLSConfig 。
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(opts *Options) {
		opts.tlsConfig = tlsConfig
	}
}

// WithIdleTimeout 。
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(opts *Options) {
		opts.idleTimeout = idleTimeout
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
	if opts.idleTimeout == 0 {
		opts.idleTimeout = defaultConnIdleTimout
	}
	if opts.dialTimeout == 0 {
		opts.dialTimeout = defaultDialTimeout
	}
	if opts.dialKeepAlive == 0 {
		opts.dialKeepAlive = defaultDialKeepAlive
	}
	return opts
}

func NewTCPConnection(endpoint string, options ...Option) (c utils.ConnReadWriteCloser, err error) {
	opts := loadOptions(options...)

	switch opts.network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, net.UnknownNetworkError(opts.network)
	}

	var conn net.Conn

	if opts.tlsConfig != nil {
		// tls connection
		dialer := tls.Dialer{
			NetDialer: &net.Dialer{
				Timeout:   opts.dialTimeout,
				KeepAlive: opts.dialKeepAlive,
				Deadline:  time.Now().Add(opts.idleTimeout),
			},
			Config: opts.tlsConfig,
		}

		conn, err = dialer.Dial(opts.network, endpoint)
		if err != nil {
			err = errors.Wrapf(err, "无法连接服务器 [endpoint=%v]", endpoint)
			return
		}
	} else {
		// normal connection
		dialer := net.Dialer{
			Timeout:   opts.dialTimeout,
			KeepAlive: opts.dialKeepAlive,
			Deadline:  time.Now().Add(opts.idleTimeout),
		}

		conn, err = dialer.Dial(opts.network, endpoint)
		if err != nil {
			err = errors.Wrapf(err, "无法连接服务器 [endpoint=%v]", endpoint)
			return
		}
	}

	c = conn.(utils.ConnReadWriteCloser)
	return
}

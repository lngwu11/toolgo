package utils

import (
	"io"
	"net"
)

type WriteCloser interface {
	io.Writer
	io.Closer
	CloseWrite() error
}

type ReadCloser interface {
	io.Reader
	io.Closer
	CloseRead() error
}

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
	CloseWrite() error
	CloseRead() error
}

type ConnReadWriteCloser interface {
	net.Conn
	CloseWrite() error
	CloseRead() error
}

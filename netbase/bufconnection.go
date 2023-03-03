package netbase

import (
	"bufio"
	"github.com/lngwu11/toolgo/utils"
)

const (
	// 默认IO流读写数据缓存大小
	defaultBufferSize = 1024 * 16
)

type bufferConnection struct {
	utils.ConnReadWriteCloser
	utils.Flusher
	rwBuffer *bufio.ReadWriter
}

func NewBufferConnection(conn utils.ConnReadWriteCloser, size int) utils.ConnReadWriteCloser {
	if size <= 0 {
		size = defaultBufferSize
	}
	buffer := bufio.NewReadWriter(
		bufio.NewReaderSize(conn, size),
		bufio.NewWriterSize(conn, size),
	)
	return &bufferConnection{
		ConnReadWriteCloser: conn,
		Flusher:             buffer,
		rwBuffer:            buffer,
	}
}

func (conn bufferConnection) Read(b []byte) (n int, err error) {
	return conn.rwBuffer.Read(b)
}

func (conn bufferConnection) Write(b []byte) (n int, err error) {
	return conn.rwBuffer.Write(b)
}

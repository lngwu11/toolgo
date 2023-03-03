package utils

type Flusher interface {
	Flush() error
}

type emptyFlusher struct{}

func (*emptyFlusher) Flush() error { return nil }

var EmptyFlusher Flusher = &emptyFlusher{}

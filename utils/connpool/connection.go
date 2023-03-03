package connpool

type IConn interface {
	lister
	Close() error
	Closed() bool
	Reusable() bool
	Detector()
}

var _ lister = (*Connection)(nil)

type lister interface {
	getNext() IConn
	getPrev() IConn
	setNext(conn IConn)
	setPrev(conn IConn)
}

type Connection struct {
	next, prev IConn
}

func (c *Connection) getNext() IConn {
	return c.next
}

func (c *Connection) getPrev() IConn {
	return c.prev
}

func (c *Connection) setNext(conn IConn) {
	c.next = conn
}

func (c *Connection) setPrev(conn IConn) {
	c.prev = conn
}
